package gatewayjwt

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var client = &http.Client{}
var gatewayID = ""
var gatewayKey = ""

func fetchAuthTokens() error {
	cmdline, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		return err
	}

	args := strings.Fields(string(cmdline))
	for _, v := range args {
		if strings.HasPrefix(v, "kg.gateway.id=") {
			gatewayID = v[14:]
		} else if strings.HasPrefix(v, "kg.gateway.key=") {
			gatewayKey = v[15:]
		}
	}

	if gatewayID == "" || gatewayKey == "" {
		return errors.New("Cannot read gateway id or key")
	}

	return nil
}

// GetGatewayJWT fetches a JWT from the IoT Identity Service used for authenticating against various Kiwigrid services.
func GetGatewayJWT(tokenURL string) (string, error) {
	if gatewayID == "" || gatewayKey == "" {
		err := fetchAuthTokens()
		if err != nil {
			return "", err
		}
	}

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		Issuer:    gatewayID,
		Subject:   gatewayID,
	}

	clientJwt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(gatewayKey))

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", tokenURL, nil)
	req.Header.Add("Authorization", "Bearer "+clientJwt)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	return data["access_token"].(string), nil
}
