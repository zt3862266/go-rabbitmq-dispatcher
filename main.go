package main

import (
	"flag"
	"fmt"
	"github.com/facebookgo/pidfile"
	"github.com/streadway/amqp"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

const (
	receiveBufferLength     = 1000
	workerBufferLength      = 1000
	receiveWorkNum          = 10
	httpMaxIdleConns        = 100
	httpMaxIdleConnsPerHost = 100
	idleConnTimeout         = 60
	consumerQosPrefetchSize = 5
	msgRetryNoLimit         = -1
	ProductionConfigDir     = "/home/rong/www/config/"
)

var usage = `Usage:[options]
Options are:
	-c configFile	Directory of the config.yaml
	-pidfile pidFile		Directory of the pidFile
	-l logFile		Directory of the logFile
`

var (
	configFile string
	logFile    string
)

type Message struct {
	Delivery   *amqp.Delivery
	SendResult bool
	QueueConf  *QueueConfig
	RetryTimes int
}

func (m *Message) SendPost(client *http.Client) int {

	client.Timeout = time.Duration(m.QueueConf.ReadTimeout) * time.Second
	return NotifyMsg(client, m.QueueConf.CallBackUrl, m.Delivery.Body)

}

func (m *Message) ack() {
	err := m.Delivery.Ack(false)
	if err != nil {
		ERROR("ack msg failed:%s", err)
	}
}

func (m *Message) nack() {
	err := m.Delivery.Nack(false, false)
	if err != nil {
		ERROR("Nack failed,err:%s", err)
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func myexit(msg string) {
	flag.Usage()
	fmt.Fprintln(os.Stderr, "\n[Error] "+msg)
	os.Exit(1)
}

func getChannel(amqpUrl string) (*amqp.Channel, error) {

	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		ERROR("get conn failed:%s", err)
		return nil, err
	}
	channel, err := conn.Channel()
	channel.Qos(consumerQosPrefetchSize, 0, false)
	if err != nil {
		ERROR("get channel failed:%s", err)
		return nil, err
	}

	return channel, nil
}

func newHttpClient(maxIdleConns, maxIdleConnsPerHost, idleConnTimeout int) *http.Client {

	transport := &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		IdleConnTimeout:     time.Duration(idleConnTimeout) * time.Second,
	}
	client := http.Client{
		Transport: transport,
	}
	return &client
}

func work(in chan Message) chan Message {

	out := make(chan Message, workerBufferLength)
	var wg sync.WaitGroup
	httpClient := newHttpClient(httpMaxIdleConns, httpMaxIdleConnsPerHost, idleConnTimeout)
	worker := func(msg Message) {

		defer wg.Done()
		if ret := msg.SendPost(httpClient); ret == CallbackResSuc {
			msg.ack()
		} else {
			msg.RetryTimes = msg.RetryTimes + 1
			if msg.QueueConf.RetryTimes == msgRetryNoLimit {
				INFO("msg retry,msg:%s,%d times", msg.Delivery.Body, msg.RetryTimes)
				in <- msg
			} else {
				if msg.RetryTimes > msg.QueueConf.RetryTimes {
					INFO("msg retry get max ,discard,retry:%d", msg.RetryTimes)
					msg.nack()
				} else {
					INFO("msg:%s, retry:%d", msg.Delivery.Body, msg.RetryTimes)
					in <- msg
				}

			}
		}

	}
	wg.Add(1)
	go func() {
		index := 0
		defer wg.Done()
		for msg := range in {
			wg.Add(1)
			go worker(msg)
			index++
			if index == 1000 {
				time.Sleep(time.Duration(1) * time.Second)
				index = 0
			}

		}
	}()

	go func() {
		wg.Wait()
		INFO("all work is done,closeing channel out")
		close(out)
	}()
	return out

}

func receive(conf *Queuesconf, envconf *EnvConfig, done chan interface{}) chan Message {

	out := make(chan Message, receiveBufferLength)
	var wg sync.WaitGroup

	receiver := func(queueConfig *QueueConfig) {
		defer wg.Done()

	reconnect:
		for {
			channel, err := getChannel(envconf.AmqpUrl)
			if err != nil {
				ERROR("get channel failed:%s,reconnect!",err)
				time.Sleep(1 * time.Second)
				continue reconnect
			}
			delivery, err := channel.Consume(
				queueConfig.Name, "", false, false, false, false, nil)
			if err != nil {
				ERROR("consumerfailed:%s,reconnect!",err)
				time.Sleep(1 * time.Second)
				continue reconnect
			}

			for {
				select {
				case msg, err := <-delivery:
					if !err {
						time.Sleep(1 * time.Second)
						INFO("channel lost,reconnect!")
						continue reconnect
					}

					outmsg := Message{
						Delivery:   &msg,
						SendResult: false,
						RetryTimes: 0,
						QueueConf:  queueConfig,
					}
					INFO("receive msg:%s", string(msg.Body))
					out <- outmsg
				case <-done:
					{
						INFO("receive done in receiver!")
						return
					}
				}

			}
		}
		INFO("recv function is finished!")

	}

	for i := 0; i < len(conf.Queues); i++ {
		myqueue := &conf.Queues[i]
		for j := 0; j < receiveWorkNum; j++ {
			wg.Add(1)
			go receiver(myqueue)
		}

	}
	go func() {
		wg.Wait()
		close(out)
		INFO("all receive gorouting is done")

	}()

	return out
}

func signalHandler(done chan interface{}) {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		s := <-c
		if s != nil {
			INFO("receive signal:,close done channel")
			close(done)
		}
	}()
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}
	flag.StringVar(&configFile, "c", "config/queue.yaml", "")
	flag.StringVar(&logFile, "l", "/home/rong/www/logs/grd.log", "")
	flag.Parse()

	if configFile == "" {
		myexit("empty configFile")
	}

	if isExist, _ := pathExists(configFile); isExist == false {
		myexit("file not exists:" + configFile)
	}
	SetlogFile(logFile)
	pidfile.Write()
	var ptrConfig *Queuesconf
	var ptrEnvConfig *EnvConfig
	if isExist, _ := pathExists(ProductionConfigDir + "queue.yaml"); isExist {
		ptrConfig = LoadAppConfig(ProductionConfigDir + "queue.yaml")
	} else {
		ptrConfig = LoadAppConfig(configFile)
	}

	if isExist, _ := pathExists(ProductionConfigDir + "config.yaml"); isExist {
		ptrEnvConfig = LoadEnvConfig(ProductionConfigDir + "config.yaml")
	} else {
		ptrEnvConfig = LoadEnvConfig("./config/config.yaml")
	}

	done := make(chan interface{}, 1)
	signalHandler(done)

	runtime.GOMAXPROCS(runtime.NumCPU())

	<-work(receive(ptrConfig, ptrEnvConfig, done))
	INFO("exit programm!")

}
