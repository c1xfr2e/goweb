package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bluecover/lm/server/codec"
	"github.com/bluecover/lm/server/render"
	"github.com/bluecover/lm/server/router"
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Start(addr string, db *gorm.DB) error {
	var (
		dec codec.Decoder
	)
	if viper.GetBool("auth.encrypt") {
		c, err := codec.NewRSACodec(viper.GetString("auth.privatekey_path"))
		if err != nil {
			return err
		}
		dec = c
		render.SetEncoder(c)
	} else {
		dec = codec.NewPlainCodec()
	}

	if viper.GetBool("debug") {
		render.PrintData()
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server := manners.NewWithServer(&http.Server{Addr: addr, Handler: router.NewRouter(db, dec)})
	go func() {
		ex := make(chan os.Signal, 1)
		signal.Notify(ex, syscall.SIGTERM, syscall.SIGINT)
		<-ex
		logrus.Info("web server shutting down...")
		server.Close()
	}()
	logrus.Infof("web server start to listen on %s", addr)
	return server.ListenAndServe()
}
