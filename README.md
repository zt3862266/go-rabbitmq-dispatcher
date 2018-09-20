# Why use it?

- 针对各类语言的consumer,开发者不需要关心消息的出错/重练/ack/reject等
- 开发者只需要编写 http restapi进行消息的消费即可
- 高性能&高并发，解决了php/python等动态语言需要启动大量消费进程的问题
- 架构层面方便运维统一管理


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

# 可用性说明

- 消费过程中出错重新连接
- 消息只有明确消费成功才会ack 并删除
- 程序异常停止,未消费的消息回到队列中


# 其他说明
- 消息通过http post body 发给到对应url, 不做任何处理
- 每个队列会默认启动10个goroutine 消费,可直接修改代码中 receiveWorkNum 修改此值
- 为了防止积压大量消费同时post 给接口造成太大压力,程序内部做了限速控制,即每post 1000条消息则sleep1秒
- 由于是并发消费,目前无法保证消息的消费顺序






