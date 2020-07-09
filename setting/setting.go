package setting

import (
	"log"
	"os"
	"strconv"
)

// S setting
var S = setting{
	AppName:       "ikuradon-salmon",
	AppVersion:    "unversioned",
	BaseURL:       "",
	APIHost:       "0.0.0.0",
	APIPort:       8080,
	UseRedis:      false,
	RedisHost:     "redis",
	RedisPort:     6379,
	RedisPassword: "",
	RedisDatabase: 0,
}

type setting struct {
	AppName       string
	AppVersion    string
	Salt          string
	BaseURL       string
	APIHost       string
	APIPort       uint16
	UseRedis      bool
	RedisHost     string
	RedisPort     uint16
	RedisPassword string
	RedisDatabase int
}

// SetSetting Add Setting
func SetSetting() {
	if appName := os.Getenv("APP_NAME"); appName != "" {
		S.AppName = appName
	}
	if appVersion := os.Getenv("APP_VERSION"); appVersion != "" {
		S.AppVersion = appVersion
	}
	if salt := os.Getenv("SALT"); salt != "" {
		S.Salt = salt
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		S.BaseURL = baseURL
	}
	if apiHost := os.Getenv("API_HOST"); apiHost != "" {
		S.APIHost = apiHost
	}
	if apiPort := os.Getenv("API_PORT"); apiPort != "" {
		apiPortuint16, err := strconv.ParseUint(os.Getenv("API_PORT"), 10, 16)
		if err != nil {
			log.Fatal(err)
		}
		S.APIPort = uint16(apiPortuint16)
	}
	if useRedis := os.Getenv("USE_REDIS"); useRedis == "true" {
		S.UseRedis = true
	}
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		S.RedisHost = redisHost
	}
	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		redisPortuint16, err := strconv.ParseUint(os.Getenv("REDIS_PORT"), 10, 16)
		if err != nil {
			log.Fatal(err)
		}
		S.RedisPort = uint16(redisPortuint16)
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		S.RedisPassword = redisPassword
	}
	if redisDatabase := os.Getenv("REDIS_DATABASE"); redisDatabase != "" {
		redisDatabaseint, err := strconv.Atoi(redisDatabase)
		if err != nil {
			log.Fatal(err)
		}
		S.RedisDatabase = redisDatabaseint
	}
}
