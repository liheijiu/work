1.kafka介绍
produccer:生产者，消息的生产者，是消息的入口

kafka cluster：kafka集群，一台多多台服务器组成

broker:是指定部署了kafka实例的服务器节点，每个服务器上一个或多
多个kafka实例，可以认为每个broker就是对应一台服务器

topic：消息的主题，可以理解为消息的分类，kafka的消息就保存在
topic，每个broker上都可以创建一个或多个topic

partition：topic的分区，每个topic可以有多个分区，分区的作用的
负载均衡，提高kafka的吞吐量，同一个topic在不同的分区的数据是不重复的
partition的表现形式就是一个一个的文件夹


replication：每个分区都有多个副本，主分区leader出现故障时候会选择一个
备胎follwer上位，成为leaderkafka中默认副本数量最大10个，，
副本的数量不能超过broker的数量，

consumer：消费者，即消息的消费方，是消息的出口


2.流程6个步骤
    1.生产者从kafka中获取leader
    2.生产者将消息发给leader
    3.leader将消息写入磁盘
    4.follwer从leader拉取数据
    5.follwer将信息写去磁盘后，给leader发送ACK
    6.leader收到所有的follower的ACK之后将把ACK发送给生产者

3.kafka选择分区的模式（3种）
    1.指定往哪个分区写
    2.指定key，kafka根据key做hash然后决定写哪个分区
    3.轮询方式

4.生产者往kafka发送数据的模式（3种）
    1. 0把数据发送给leader就成功
    2. 把数据发送给leader，等leader回复ACK就成功
    3.把数据发送给leader，follower从leader拉取数据，
    回复leaderACK，leader将ACK在回复生产者
    
logagent：    

1.logagent工作流程：
    1读日志----tailf
    2往kafka写日志------sarama
    













