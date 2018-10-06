package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bluecover/lm/command"
	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/module/crm"
	"github.com/bluecover/lm/server"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zaoshu/hardcore/logging"
	"github.com/zaoshu/hermogo"
)

func initConfig() {
	viper.SetConfigName("default")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("LEMONBE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func initDB() *gorm.DB {
	driver := viper.GetString("db.driver")
	dsn := viper.GetString("db.dsn")
	db, err := gorm.Open(driver, dsn)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(models.GetModels()...).Error
	if err != nil {
		panic(err)
	}
	return db
}

func initHermogo() {
	hermogoConfig := hermogo.Config{
		Endpoint:  os.Getenv("HermogoEndpoint"),
		AccessID:  os.Getenv("HermogoAccessID"),
		AccessKey: os.Getenv("HermogoAccessKey"),
	}

	logrus.Info("hermogo", os.Getenv("HermogoEndpoint"))
	if err := hermogo.Init(hermogoConfig); err != nil {
		panic(err)
	}
}

func initCRM() {
	config := crm.Xiaoshouyi{}
	err := viper.UnmarshalKey("xiaoshouyi", &config)
	if err != nil {
		panic(err)
	}
	crm.Init(config)
}

func main() {
	initConfig()

	logging.InitFromConfig(logging.Config{
		Level:     viper.GetString("log.level"),
		Formatter: viper.GetString("log.formatter"),
	})

	db := initDB()
	defer db.Close()

	if len(os.Args) > 1 {
		command.RunCommand(db)
		return
	}

	if !viper.GetBool("develop.disable_mns") {
		initHermogo()
	}

	initCRM()

	addr := fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt("server.port"))
	server.Start(addr, db)
}
