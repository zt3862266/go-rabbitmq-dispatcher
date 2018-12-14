package main

import (
	"testing"
	"github.com/streadway/amqp"
)

func TestPub(t *testing.T){

	connection,err := amqp.Dial("amqp://rdguest:Grong@online@10.12.33.239:5672/")
	if err != nil {
		t.Error("get conneciton failed",err)
	}
	defer connection.Close()

	channel,err := connection.Channel()
	if err != nil{
		t.Error("get channel failed",err)

	}

	for i:=0;i<10000;i++ {
		if err = channel.Publish(
				"",
				"mq-test",
				false,
				false,
				amqp.Publishing{
					Headers:amqp.Table{},
					ContentType:"text/plain",
					ContentEncoding:"",
					Body: []byte("你好,世界人民大和平!"),
					DeliveryMode:amqp.Persistent,
					Priority:0,
				},
			);
		err !=nil{
			t.Error("publish failed")
		}

	}
}
