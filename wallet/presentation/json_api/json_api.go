package jsonapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ageeknamedslickback/wallet-API/wallet/dto"
	"github.com/ageeknamedslickback/wallet-API/wallet/usecases"
	"github.com/gin-gonic/gin"
)

type WalletJsonPresentation interface {
	Authenticate(c *gin.Context)

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

	var crAmountInput dto.AmountInput
	if err := c.ShouldBindJSON(&crAmountInput); err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := crAmountInput.Valid(); err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	wallet, err := p.Uc.CreditWallet(
		ctx,
		*walletID,
		crAmountInput.Amount,
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

	var drAmountInput dto.AmountInput
	if err := c.ShouldBindJSON(&drAmountInput); err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := drAmountInput.Valid(); err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	wallet, err := p.Uc.DebitWallet(
		ctx,
		*walletID,
		drAmountInput.Amount,
	)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}

// Authenticate provides an authentication endpoint that returns an access token
// to interact with the other APIs
func (p *WalletJsonAPI) Authenticate(c *gin.Context) {
	params := url.Values{}
	params.Add("grant_type", os.Getenv("AUTH0_GRANT_TYPE"))
	params.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	params.Add("client_secret", os.Getenv("AUTH0_CLIENT_SECRET"))
	params.Add("audience", os.Getenv("AUTH0_AUDIENCE"))
	payload := strings.NewReader(params.Encode())

	URL := fmt.Sprintf("https://%s/oauth/token", os.Getenv("AUTH0_DOMAIN"))
	req, err := http.NewRequest(http.MethodPost, URL, payload)
	if err != nil {
		jsonErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		jsonErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		jsonErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var accessToken dto.AccessToken
	if err := json.Unmarshal(body, &accessToken); err != nil {
		jsonErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": accessToken})
}
