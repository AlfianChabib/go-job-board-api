package config

import (
	"github.com/spf13/viper"
)

type Env struct {
	AppEnv            string
	Port              string
	DatabaseUser      string
	DatabasePassword  string
	DatabaseDB        string
	DatabaseHost      string
	DatabasePort      int
	DatabaseUrl       string
	AccessSecretKey   string
	AccessDuration    int
	RefreshSecretKey  string
	RefreshDuration   int
	MinioRootUser     string
	MinioRootPassword string
	MinioAvatarBucket string
	MinioCvBucket     string
	MinioUseSSL       bool
	MinioEndpoint     string
	MinioPublicUrl    string
}

func LoadEnv() *Env {
	config := viper.New()
	config.SetConfigFile(".env")

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var env *Env = &Env{
		AppEnv:            config.GetString("APP_ENV"),
		Port:              config.GetString("PORT"),
		DatabaseUser:      config.GetString("DATABASE_USER"),
		DatabasePassword:  config.GetString("DATABASE_PASSWORD"),
		DatabaseDB:        config.GetString("DATABASE_DB"),
		DatabaseHost:      config.GetString("DATABASE_HOST"),
		DatabasePort:      config.GetInt("DATABASE_PORT"),
		DatabaseUrl:       config.GetString("DATABASE_URL"),
		AccessSecretKey:   config.GetString("ACCESS_SECRET_KEY"),
		AccessDuration:    config.GetInt("ACCESS_DURATION"), // second
		RefreshSecretKey:  config.GetString("REFRESH_SECRET_KEY"),
		RefreshDuration:   config.GetInt("REFRESH_DURATION"), // second
		MinioRootUser:     config.GetString("MINIO_ROOT_USER"),
		MinioRootPassword: config.GetString("MINIO_ROOT_PASSWORD"),
		MinioAvatarBucket: config.GetString("MINIO_AVATAR_BUCKET"),
		MinioCvBucket:     config.GetString("MINIO_CV_BUCKET"),
		MinioUseSSL:       config.GetBool("MINIO_USE_SSL"),
		MinioEndpoint:     config.GetString("MINIO_ENDPOINT"),
		MinioPublicUrl:    config.GetString("MINIO_PUBLIC_URL"),
	}

	return env
}
