# go-rabbitmq-dispatcher

A simple rabbitmq message dispatcher writen in golang

# Why user it?

- 针对各类语言的consumser,开发者不需要关心消息的出错/重练/ack/reject等
- 开发者只需要编写 http restapi进行消息的消费即可
- 高性能&高并发，解决了php/python等动态语言需要启动大量消费进程的问题

# how to use it ?

安装 golang

下载代码至 /PATH/to/go-rabbitmq-dispatcher

go build

./start.sh

# 配置文件说明

config/queue.yaml:队列的配置，可配置多个队列的信息，详细配置项如下：

- name:队列名称
- call_back_url:该队列消费者对应的url
- read_timeout:请求超时时间，单位秒
- retry_times: 请求失败重试次数，若设置为-1，则无限次重试,直到成功

config/config.yaml: 连接队列的uri 格式为:

- amqp://user:password@ip:port/vhost
 

# call_back_url 返回格式：

配置文件中配置的call_back_url务必以以下格式返回,这样方可判断消息是否消费成功or失败

消息处理成功:

{
    "error": 0, "msg": ""
}

消息处理失败:

{
    "error": 1,"msg": "time out" 
}

error: 0 正常 1 异常

msg: 返回的消息内容
 


