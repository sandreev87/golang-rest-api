package app

import (
	"github.com/coocood/freecache"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sandreev87/golang-rest-api/internal/handler"
	"github.com/sandreev87/golang-rest-api/internal/repository"
	"github.com/sandreev87/golang-rest-api/internal/server"
	"github.com/sandreev87/golang-rest-api/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

const CacheSize = 104857600 // 100MB

func Run() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	refreshTokenCache := freecache.NewCache(CacheSize)

	repos := repository.NewRepository(initDbConnection(), refreshTokenCache)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
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

func initDbConnection() *sqlx.DB {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	return db
}
