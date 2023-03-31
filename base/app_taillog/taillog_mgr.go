package app_taillog

import (
	"fmt"
	"time"
	"work/base/app_etcd"
)

var tskMgr *tailiMgr

//TailTask管理者
type tailiMgr struct {
	etcdConfkafka []*app_etcd.EtcdConfKafa
	tskMap        map[string]*TailTask
	newConfChan   chan []*app_etcd.EtcdConfKafa
}

func Init(etcdConfKafka []*app_etcd.EtcdConfKafa) {
	tskMgr = &tailiMgr{
		etcdConfkafka: etcdConfKafka, //把日志收集项配置信息保存起来
		tskMap:        make(map[string]*TailTask, 16),
		newConfChan:   make(chan []*app_etcd.EtcdConfKafa), //无缓冲区通道
	}
	//循环每一个日志收集项，创建topic
	for _, etcdConf := range etcdConfKafka {

		tailObj := NewTailTask(etcdConf.Path, etcdConf.Topic) //初始化的时候起了多少都需要记下来，为了后续判断方便
		mk := fmt.Sprintf("%s_%s", etcdConf.Path, etcdConf.Topic)
		tskMgr.tskMap[mk] = tailObj
	}
	go tskMgr.run()

}

//监听自己的newConfChan， 有了新的配置更新后多对应的操作
func (t *tailiMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.tskMap[conf.Path]
				if ok {
					//原来就有,不需要操作
					continue
				} else {
					//新增
					tailObj := NewTailTask(conf.Path, conf.Topic)
					t.tskMap[mk] = tailObj
				}
			}
			for _, c1 := range t.etcdConfkafka {
				isDelete := true
				for _, c2 := range newConf {
					if c2.Path == c1.Path && c2.Topic == c1.Topic {
						isDelete = false
						continue
					}
				}
				if isDelete {
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.tskMap[mk].cancelFunc()
				}
			}
			fmt.Println("新的配置来了", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

//向外暴露一个函数,项 tskMgr 的newConfChan
func NewConf() chan<- []*app_etcd.EtcdConfKafa {
	return tskMgr.newConfChan
}
