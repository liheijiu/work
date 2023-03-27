package runner

import (
	"fmt"
	"gopkg.in/ini.v1"
	"time"
	"work/base/app_conf"
	"work/base/app_kafka"
	"work/base/app_mysql"
	"work/base/app_taillog"
)

var cfg = new(app_conf.AppConf)

func Run() {
	//初始化配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("ini.MapTo failed, err:%v\n", err)
		return
	}
	fmt.Println("Init Conf success...")

	//初始化 mysql 连接
	err = app_mysql.InitDb(cfg)
	if err != nil {
		fmt.Printf("init  mysql failed, err:%v\n", err)
		return
	}
	fmt.Println("The connection to the database was successful")

	//kafka 初始化
	err = app_kafka.InitKafka(cfg)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	fmt.Println("init kafka success....")

	//tailog  初始化
	err = app_taillog.InitTaillog(cfg.TaillogConf.FileName)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	fmt.Println("Init taillog success")

	//单条查询
	//app_mysql.MysqlQuire(9)
	//多条查询
	//app_mysql.MysqlQuires(2)

	//插入数据
	//app_mysql.MysqInsert()

	//更新数据
	//app_mysql.MysqUpdate(100, 5)

	//删除数据
	//app_mysql.MysqDelete(5)

	//sql预处理
	//app_mysql.PrepareInsertDemo()

	//事务操作
	//app_mysql.TransactionDemo()

	//往kafka发送消息
	for {
		select {
		//读取的日志添加到chan中
		case line := <-app_taillog.ReadChan():

			//从通道中拿取日志，写入kafka指定的topic
			app_kafka.SendTokafka(cfg.KafkaConf.Topic, line.Text)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}
