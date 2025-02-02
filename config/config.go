package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App        App
		Gcp        Gcp
		QrcodeInfo QrcodeInfo
		User       User
	}

	App struct {
		Name   string
		Token  string
		AppUrl string
	}

	QrcodeInfo struct {
		AccountInfo string
	}

	Gcp struct {
		BucketName string
		UploadPath string
		FileLimit  int64
	}

	User struct {
		Username string
		Password string
	}
)

func GetConfig(path string) Config {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(fmt.Sprintf("./" + path))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("fatal error read config filed %v", err.Error())
		panic(err)
	}

	return Config{
		App: App{
			Name:   viper.GetString("APP_NAME"),
			Token:  viper.GetString("APP_TOKEN"),
			AppUrl: viper.GetString("APP_URL"),
		},
		QrcodeInfo: QrcodeInfo{
			AccountInfo: viper.GetString("PAYMENT_INFO"),
		},
		Gcp: Gcp{
			BucketName: viper.GetString("GCP_BUCKET_NAME"),
			UploadPath: viper.GetString("GCP_UPLOAD_PATH"),
			FileLimit:  int64(viper.GetInt("GCP_UPLOAD_PATH")),
		},
		User: User{
			Username: viper.GetString("IG_USERNAME"),
			Password: viper.GetString("IG_PASSWORD"),
		},
	}
}
