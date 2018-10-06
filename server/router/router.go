package router

import (
	"github.com/bluecover/lm/server/codec"
	"github.com/bluecover/lm/server/handler"
	"github.com/bluecover/lm/server/middleware"
	"github.com/bluecover/lm/server/middleware/authware"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zaoshu/hardcore/logging/ginware"
	"github.com/zaoshu/hardcore/trace/gintrace"
)

func NewRouter(db *gorm.DB, dec codec.Decoder) *gin.Engine {
	router := gin.New()
	router.HandleMethodNotAllowed = true
	router.Use(gintrace.WithRequestID(true), ginware.Logger())
	router.Use(middleware.RequestDecoder(dec))

	addREST(router, db)

	return router
}

func addREST(r *gin.Engine, db *gorm.DB) {

	// Authentication API
	r.POST("/REST/user/signin", handler.Login(db))
	r.POST("/REST/user/renewToken", handler.RenewToken(db))
	r.POST("/REST/user/requestTrial", handler.RequestTrial(db))

	// Authenticated group
	authGroup := r.Group("/REST/", authware.LoggedIn(db))
	authGroup.POST("ping", handler.Ping)

	authGroup.GET("user/info", handler.UserInfo(db))
	authGroup.GET("dataset/list", handler.DatasetList(db))
	authGroup.GET("dataset/figurePage", handler.GetFingerPage(db))
	authGroup.GET("dataset/figure", handler.GetFinger(db))
	authGroup.GET("dataset/export", handler.DataExport(db))
	authGroup.GET("dataset/filter", handler.Filter(db))

	// search
	searchOp := &handler.SearchOp{DB: db}
	searchOp.LoadAllFigures()
	authGroup.GET("dataset/search", searchOp.Search())
	authGroup.POST("dataset/search", searchOp.SaveSearchHistory())

	// add request
	requestOp := &handler.RequestOp{DB: db}
	authGroup.POST("dataset/addRequest", requestOp.NewRequest())
}
