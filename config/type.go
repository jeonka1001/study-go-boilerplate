package config

type Config struct {
	Server Server `mapstructure:"server"`
	App    App    `mapstructure:"app"`
	DB     DB     `mapstructure:"db"`
}

type Server struct {
	Profile string `mapstructure:"profile"`
	Port    string `mapstructure:"port"`
}

type App struct {
	Version string `mapstructure:"version"`
	Url     Url    `mapstructure:"url"`
	Log     Log    `mapstructure:"log"`
}

// Url: 이 서비스가 호출하는 외부 API base-url (필요에 따라 추가/변경)
type Url struct {
	Account struct {
		BaseUrl string `mapstructure:"base-url"`
	} `mapstructure:"account"`
}

type Log struct {
	Level       string `mapstructure:"level"`
	HistoryType string `mapstructure:"history-type"`
}

// DB: MySQL/MariaDB 접속 정보. Password 는 yml 에 실제 값을 직접 넣지 말고
// 환경변수(DB_PASSWORD)로 오버라이드하거나 배포 시크릿으로 주입할 것.
type DB struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	Charset         string `mapstructure:"charset"`
	MaxOpenConns    int    `mapstructure:"max-open-conns"`
	MaxIdleConns    int    `mapstructure:"max-idle-conns"`
	ConnMaxLifetime int    `mapstructure:"conn-max-lifetime-minutes"`
}
