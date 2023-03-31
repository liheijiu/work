package app_conf

//kafka  config
type KafkaConf struct {
	Address     string `ini:"address"`
	ChanMaxSize int    `ini:"chanmaxsize"`
}

//etcd  config
type EtcdConf struct {
	Address string `ini:"address"`
	Key     string `ini:"key"`
	Timeout int    `ini:"timeout"`
}

//kafka  , etcd config
type AppConf struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
}

//
