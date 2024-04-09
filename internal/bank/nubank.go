package bank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/mdp/qrterminal/v3"
	"github.com/reonardoleis/manager/internal/models"
	"github.com/reonardoleis/manager/internal/service"
)

const DISCOVERY_URL = "https://prod-global-webapp-proxy.nubank.com.br/api/discovery"
const DISCOVERY_APP_URL = "https://prod-global-webapp-proxy.nubank.com.br/api/app/discovery"

var accessToken string
var revokeUrl string

type Nubank struct {
	cpf          string
	password     string
	clientSecret string
}

func NewNubank(cpf, password, clientSecret string) service.Bank {
	return Nubank{cpf, password, clientSecret}
}

func (n Nubank) GetBill(url string) (*models.Bill, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	v := new(bytes.Buffer)
	v.ReadFrom(response.Body)

	m := make(map[string]interface{})
	err = json.Unmarshal(v.Bytes(), &m)
	if err != nil {
		return nil, err
	}

	bill := m["bill"].(map[string]interface{})

	remarshalled, err := json.Marshal(bill["line_items"])
	if err != nil {
		return nil, err
	}

	txs := []models.Tx{}

	err = json.Unmarshal(remarshalled, &txs)
	if err != nil {
		return nil, err
	}

	b := new(models.Bill)
	b.Txs = txs
	return b, nil
}

func (n Nubank) Authorize() error {
	discovery := new(discoveryResponse)

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", DISCOVERY_URL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(discovery)
	if err != nil {
		return err
	}

	loginUrl := discovery.Login

	loginBody, err := json.Marshal(authorizeRequest{
		ClientID:     "other.conta",
		ClientSecret: n.clientSecret,
		GrantType:    "password",
		Login:        n.cpf,
		Password:     n.password,
	})
	if err != nil {
		return err
	}

	req, err = http.NewRequest("POST", loginUrl, bytes.NewBuffer(loginBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	response, err = httpClient.Do(req)
	if err != nil {
		return err
	}

	login := new(authorizeResponse)
	err = json.NewDecoder(response.Body).Decode(login)
	if err != nil {
		return err
	}

	httpClient = &http.Client{}
	req, err = http.NewRequest("GET", DISCOVERY_APP_URL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	response, err = httpClient.Do(req)
	if err != nil {
		return err
	}

	discoveryApp := new(discoveryAppResponse)
	err = json.NewDecoder(response.Body).Decode(discoveryApp)
	if err != nil {
		return err
	}

	qrCodeId := uuid.New().String()
	qrterminal.Generate(qrCodeId, qrterminal.H, os.Stdout)

	liftBody, err := json.Marshal(liftRequest{
		QrCodeId: qrCodeId,
		Type:     "login-webapp",
	})
	if err != nil {
		return err
	}

	fmt.Println("Scan the QR code above and press enter to continue")
	var answer string
	fmt.Scanln(&answer)

	req, err = http.NewRequest("POST", discoveryApp.Lift, bytes.NewBuffer(liftBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+login.AccessToken)

	response, err = httpClient.Do(req)
	if err != nil {
		return err
	}

	lift := new(liftResponse)
	err = json.NewDecoder(response.Body).Decode(lift)
	if err != nil {
		return err
	}

	accessToken = lift.AccessToken
	revokeUrl = lift.Links.RevokeToken.Href

	return nil
}

func (n Nubank) Revoke() error {
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", revokeUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	_, err = httpClient.Do(req)
	if err != nil {
		log.Println("error revoking token", err)
		return err
	}

	log.Println("token revoked!")

	return nil
}
