package presentation

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/ageeknamedslickback/wallet-API/wallet/infrastructure/database"
	jsonapi "github.com/ageeknamedslickback/wallet-API/wallet/presentation/json_api"
	"github.com/ageeknamedslickback/wallet-API/wallet/presentation/middleware"
	"github.com/ageeknamedslickback/wallet-API/wallet/usecases"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
)

// Router sets up the presentation layer config router
func Router() *gin.Engine {
	router := gin.Default()

	gormDb, err := database.ConnectToDatabase()
	if err != nil {
		log.Panicf("error connecting to the database: %v", err)
	}
	getRepo := database.NewWalletDb(gormDb)
	updateRepo := database.NewWalletDb(gormDb)
	uc := usecases.NewWalletUsecases(getRepo, updateRepo)
	h := jsonapi.NewWalletJsonAPIs(uc)

	gin.DisableConsoleColor()

	f, _ := os.Create("wallet.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.POST("/access_token", h.Authenticate)

	v1 := router.Group("api/v1")
	v1.Use(adapter.Wrap(middleware.EnsureValidToken()))
	{
		v1.GET("/:wallet_id/balance", h.WalletBalance)
		v1.POST("/:wallet_id/credit", h.CreditWallet)
		v1.POST("/:wallet_id/debit", h.DebitWallet)
	}

	return router
}
