package runner

import (
	"fmt"
	"gopkg.in/ini.v1"
	"sync"
	"time"
	"work/base/app_taillog"

	//"time"
	"work/base/app_conf"
	"work/base/app_etcd"
	"work/base/app_kafka"
	//"work/base/app_taillog"
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

	//kafka 初始化
	err = app_kafka.InitKafka(cfg.KafkaConf.Address, cfg.KafkaConf.ChanMaxSize)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	fmt.Println("init kafka success....")

	//etcd 初始化
	err = app_etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("init etcd success....")

	//从etcd的key拉取kafka的配置项信息
	etcdConfKafka, err := app_etcd.GetEtcdConf(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Printf("from etcd get config failed,err %v\n", err)
		return
	}

	//收集的日志发往kafka
	app_taillog.Init(etcdConfKafka)

	newConf := app_taillog.NewConf() //获取对外暴露的通道
	var wg sync.WaitGroup
	wg.Add(1)
	go app_etcd.WatchConf(cfg.EtcdConf.Key, newConf) //获取最新更新配置通知上面的通道
	wg.Wait()
}
