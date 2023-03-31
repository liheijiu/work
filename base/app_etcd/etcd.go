package app_etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	client *clientv3.Client
)

//etcd 的key中存放kafka配置的信息项
type EtcdConfKafa struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

//etcd  做初始化->
func Init(address string, timeout time.Duration) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{address},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	return
}

//根据etcd 的key 获取配置项
func GetEtcdConf(key string) (etcdConfKafka []*EtcdConfKafa, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, key)
	cancel()
	if err != nil {
		return
	}
	//遍历键值对
	for _, kv := range resp.Kvs {
		err = json.Unmarshal(kv.Value, &etcdConfKafka)
		if err != nil {
			fmt.Printf("Unmarshal  etcd value failed,err %v\n", err)
			return
		}
	}
	//遍历 kafka的配置信息
	fmt.Printf("get etcd.conf success, %v\n", etcdConfKafka)
	for index, value := range etcdConfKafka {
		fmt.Printf("index:%v  value:%v\n", index, value)
	}
	return
}

//etcd  watch
func WatchConf(key string, newConfCh chan<- []*EtcdConfKafa) {
	ch := client.Watch(context.Background(), key) //监控etcd中key的变化-创建、更改、删除
	//从通道尝试获取值（监视的信息
	for wresp := range ch {
		for _, evn := range wresp.Events {
			fmt.Printf("type:%v  key:%v  value:%v\n", evn.Type, string(evn.Kv.Key), string(evn.Kv.Value))

			//通知别人
			//先判断类型
			var newConf []*EtcdConfKafa
			if evn.Type != clientv3.EventTypeDelete {
				//如果删除操作,手动传递一个空的配置
				err := json.Unmarshal(evn.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("json.Unmarshal failed, err:%v\n", err)
					continue
				}
			}
			fmt.Printf("get new conf:%v\n", newConf)
			newConfCh <- newConf
		}
	}
}
