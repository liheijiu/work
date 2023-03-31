package app_taillog

import (
	"context"
	"github.com/hpcloud/tail"
	"time"
	"work/base/app_kafka"
)

//TailTask 日志收集的实例
type TailTask struct {
	path       string             //路径
	topic      string             //topic
	instance   *tail.Tail         //实例
	ctx        context.Context    //为了能够实现退出t.run
	cancelFunc context.CancelFunc //为了能够实现退出t.run
}

//采集日志
func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background()) //当不存在配置时，退出G使用
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	tailObj.init() //根据路径去打开对应的日志
	return
}

//TailTask  初始化
func (t *TailTask) init() {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,                                 //
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		return
	}
	go t.run() //直接去采集日志发送kafka数据
}

//读取日志
func (t *TailTask) run() {
	for {
		select {
		case line := <-t.instance.Lines: //从instance通道一行一行的读取日志

			app_kafka.SendToChan(t.topic, line.Text) //日志发送到一个通道中

		default:

			time.Sleep(time.Millisecond * 500) //如果没有读取到就休息
		}
	}
}
