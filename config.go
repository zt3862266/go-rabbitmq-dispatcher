package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Queuesconf struct {
	Queues []QueueConfig `yaml:"queues"`
}

type QueueConfig struct {
	Name        string `yaml:"name"`
	CallBackUrl string `yaml:"call_back_url"`
	ReadTimeout int    `yaml:"read_timeout"`
	RetryTimes  int    `yaml:"retry_times"`
}
type EnvConfig struct {
	AmqpUrl string `yaml:"amqp_url"`
}

func LoadAppConfig(configFile string) *Queuesconf {

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	myconf := &Queuesconf{}
	err = yaml.Unmarshal(file, myconf)
	if err != nil {
		panic(err)
	}
	return myconf
}

func LoadEnvConfig(configFile string) *EnvConfig {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	myconf := &EnvConfig{}
	err = yaml.Unmarshal(file, myconf)
	if err != nil {
		panic(err)
	}
	return myconf
}
