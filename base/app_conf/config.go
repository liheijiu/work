package app_conf

//mysql  config
type AppMysql struct {
	Address  string `ini:"address"`
	Database string `ini:"database"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

//redis  config

//kafka  config
type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

//tailflog   config
type TaillogConf struct {
	FileName string `ini:"fileName"`
}

//kafka  , taillog , mysql config
type AppConf struct {
	KafkaConf   `ini:"kafka"`
	TaillogConf `ini:"taillog"`
	AppMysql    `ini:"mysql"`
}

//
