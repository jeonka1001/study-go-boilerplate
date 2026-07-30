package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
)

// Profile 은 실행 환경을 나타낸다. (local/dev/qa/prod)
var Profile string

// envKeyReplacer: viper 키의 "." 을 환경변수 규칙인 "_" 로 치환 (예: db.password -> DB_PASSWORD)
var envKeyReplacer = strings.NewReplacer(".", "_")

func init() {
	Profile = os.Getenv("APP_ENV")
	if len(Profile) == 0 {
		Profile = "local"
	}
}

// NewConfig 는 app/config/yml/<profile>.yml 을 읽어 Config 로 파싱한다.
// DB 비밀번호 등 민감 정보는 yml 에 직접 넣지 말고 환경변수로 오버라이드하는 것을 권장한다.
func NewConfig() (*Config, error) {
	var cfg *Config
	path := filepath.Join(fmt.Sprintf("app/config/yml/%s.yml", Profile))

	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)

	// yml 에 없는 값은 환경변수로 주입 가능 (예: DB_PASSWORD -> db.password)
	viper.SetEnvKeyReplacer(envKeyReplacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Cannot read %s: %v", path, err)
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
