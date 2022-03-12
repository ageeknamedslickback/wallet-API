package presentation

import (
	"log"

	"github.com/ageeknamedslickback/wallet-API/wallet/infrastructure/database"
	jsonapi "github.com/ageeknamedslickback/wallet-API/wallet/presentation/json_api"
	"github.com/ageeknamedslickback/wallet-API/wallet/usecases"
	"github.com/gin-gonic/gin"
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

	v1 := router.Group("api/v1")
	{
		v1.GET("/:wallet_id/balance", h.WalletBalance)
		v1.POST("/:wallet_id/credit", h.CreditWallet)
		v1.POST("/:wallet_id/debit", h.DebitWallet)
	}

	return router
}
