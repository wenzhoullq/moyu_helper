package config

type Configuration struct {
	ServerConfigure   `toml:"server_configure"`
	QianFanConfigure  `toml:"qianfan_configure"`
	MysqlConfigure    `toml:"mysql_configure"`
	LogConfigure      `toml:"log_configure"`
	RedisConfigure    `toml:"redis_configure"`
	SignConfigure     `toml:"sign_configure"`
	FileConfigure     `toml:"file_configure"`
	CronTaskConfigure `toml:"corn_task_configure"`
	TanshuConfigure   `toml:"tanshu_configure"`
	TxConfigure       `toml:"tx_configure"`
}

type TxConfigure struct {
	SecretId  string `toml:"secret_id"`
	SecretKey string `toml:"secret_key"`
}

type FileConfigure struct {
	HolidayFile string `toml:"holiday_file"`
	ImgFile     string `toml:"img_file"`
}

type SignConfigure struct {
	SignRewardFirst  int    `toml:"sign_reward_first"`
	SignRewardSecond int    `toml:"sign_reward_second"`
	SignRewardThird  int    `toml:"sign_reward_third"`
	SignRewardElse   int    `toml:"sign_reward_else"`
	SignMark         string `toml:"sign_mark"`
}

type RedisConfigure struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type LogConfigure struct {
	LogFilePath string `toml:"log_file_path"`
	LogFileName string `toml:"log_file_name"`
}

type ServerConfigure struct {
	ServerAddr string `toml:"server_addr"`
}
type MysqlConfigure struct {
	Driver   string `toml:"driver"`
	UserName string `toml:"user_name"`
	Pw       string `toml:"pw"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	DbName   string `toml:"dbname"`
	TimeOut  string `toml:"timeout"`
}

type QianFanConfigure struct {
	ApiKey    string `toml:"api_key"`
	SecretKey string `toml:"secret_key"`
	AccessKey string `toml:"access_key"`
}

type TanshuConfigure struct {
	Key string `toml:"key"`
}

type CronTaskConfigure struct {
	HolidayTips   string `toml:"holiday_tips"`
	NewsTips      string `toml:"news_tips"`
	RegularUpdate string `toml:"regular_update"`
}
