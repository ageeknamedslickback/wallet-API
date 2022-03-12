package jsonapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ageeknamedslickback/wallet-API/wallet/dto"
	"github.com/ageeknamedslickback/wallet-API/wallet/usecases"
	"github.com/gin-gonic/gin"
)

type WalletJsonPresentation interface {
	WalletBalance(c *gin.Context)
	CreditWallet(c *gin.Context)
	DebitWallet(c *gin.Context)
}

// WalletJsonAPI sets up wallet's API server presentation layer
// with all the necessary dependencies
type WalletJsonAPI struct {
	Uc usecases.WalletBusinessLogic
}

// NewWalletJsonAPIs initializes a new instance of wallet's JSON APIs
func NewWalletJsonAPIs(uc usecases.WalletBusinessLogic) *WalletJsonAPI {
	w := &WalletJsonAPI{
		Uc: uc,
	}
	w.checkPreconditions()
	return w
}

func (p *WalletJsonAPI) checkPreconditions() {
	if p.Uc == nil {
		log.Panicf("presentation layer has not initialized the usecases")
	}
}

func jsonErrorResponse(c *gin.Context, statusCode int, err string) {
	c.JSON(statusCode, gin.H{"error": err})
}

func getWalletID(c *gin.Context) (*int, error) {
	strWalletID := c.Param("wallet_id")
	if strWalletID == "" {
		return nil, fmt.Errorf("wallet ID has not been provided")
	}

	walletID, err := strconv.Atoi(strWalletID)
	if err != nil {
		return nil, err
	}

	return &walletID, nil
}

// WalletBalance is a JSON API that retrieves a wallet's balance
func (p *WalletJsonAPI) WalletBalance(c *gin.Context) {
	ctx := context.Background()

	walletID, err := getWalletID(c)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	balance, err := p.Uc.WalletBalance(ctx, *walletID)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// CreditWallet is a JSON API that credits a wallet's balance
func (p *WalletJsonAPI) CreditWallet(c *gin.Context) {
	ctx := context.Background()

	walletID, err := getWalletID(c)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var crAmountInput dto.CrAmountInput
	if err := c.ShouldBindJSON(&crAmountInput); err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	wallet, err := p.Uc.CreditWallet(
		ctx,
		*walletID,
		crAmountInput.CreditAmount,
	)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}

// DebitWallet is a JSON API that credits a wallet's balance
func (p *WalletJsonAPI) DebitWallet(c *gin.Context) {
	ctx := context.Background()

	walletID, err := getWalletID(c)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var drAmountInput dto.DrAmountInput
	if err := c.ShouldBindJSON(&drAmountInput); err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	wallet, err := p.Uc.DebitWallet(
		ctx,
		*walletID,
		drAmountInput.DebitAmount,
	)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}
