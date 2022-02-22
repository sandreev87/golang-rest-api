package app

import (
	server "github.com/sandreev87/golang-rest-api"
	"github.com/sandreev87/golang-rest-api/internal/handler"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing configs: %s", err.Error())
	}

	handlers := handler.NewHandler()
	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
