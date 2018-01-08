package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/net/context"
)

const (
	TOKEN_URL = "https://api.%s.onelogin.com/auth/oauth2/v2/token"
	AUTH_URL = "https://api.%s.onelogin.com/auth/oauth2/auth"
	BASE_PATH = "https://api.%s.onelogin.com/api/1/%s"
)

type OneloginClient struct {
	*http.Client
	region string
}

func NewOneloginClient(clientId, clientSecret, region string) (*OneloginClient, error) {
	ctx := context.Background()
	conf := &OAuth2Config{&oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"saml_assertion"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf(AUTH_URL, region),
			TokenURL: fmt.Sprintf(TOKEN_URL, region),
		},
	}}

	token, err := conf.GetToken()
	if err != nil {
		return nil, err
	}

	client := conf.GetClient(ctx, token, region)
	return client, nil
}

func (c *OneloginClient) PostJson(url string, body interface{}) (*http.Response, error) {
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.Post(url, "application/json", bytes.NewReader(json))
}

func (c *OneloginClient) GetSamlAssertion(appId, subdomain, username, password, mfaCode string) (*SamlAssertion, error) {
	req_params := &samlAssertionRequest{
		UsernameOrEmail: username,
		Password: password,
		AppId: appId,
		Subdomain: subdomain,
	}
	resp, err := c.PostJson(fmt.Sprintf(BASE_PATH, c.region, "saml_assertion"), req_params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	samlResponse := &samlAssertionResponse{}
	err = parseJson(resp, samlResponse)
	if err != nil {
		return nil, err
	}

	// TODO fix workflow w/o mfa

	verify_req_params := &verifyFactorRequest{
		AppId: appId,
		DeviceId: samlResponse.Data[0].Devices[0].DeviceId,
		StateToken: samlResponse.Data[0].StateToken,
		OtpToken: mfaCode,
	}
	resp, err = c.PostJson(samlResponse.Data[0].CallbackUrl, verify_req_params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	verifyResponse := &verifyFactorResponse{}
	err = parseJson(resp, verifyResponse)
	if err != nil {
		return nil, err
	}

	assertion := NewSamlAssertion(verifyResponse.Data)
	return assertion, nil
}

type OAuth2Config struct {
	*oauth2.Config
}

func (conf *OAuth2Config) GetToken() (*oauth2.Token, error) {
	req, err := http.NewRequest("POST", conf.Endpoint.TokenURL, strings.NewReader(`{"grant_type":"client_credentials"}`))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("client_id:%s, client_secret:%s", conf.ClientID, conf.ClientSecret))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	json.Unmarshal(body, token)

	return token, nil
}

func (conf *OAuth2Config) GetClient(ctx context.Context, t *oauth2.Token, region string) *OneloginClient {
	return &OneloginClient{conf.Client(ctx, t), region}
}

type samlAssertionRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
	AppId           string `json:"app_id"`
	Subdomain       string `json:"subdomain"`
	IpAddress       string `json:"ip_address,omitempty"`
}

type samlAssertionResponseDataDeveice struct {
	DeviceId   int `json:"device_id"`
	DeviceType string `json:"device_type"`
}

type samlAssertionResponseData struct {
	StateToken  string `json:"state_token"`
	Devices     []*samlAssertionResponseDataDeveice `json:"devices"`
	CallbackUrl string `json:"callback_url"`
}

type responseStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int `json:"code"`
	Error   bool `json:"error"`
}

type samlAssertionResponse struct {
	Status *responseStatus `json:"status"`
	Data   []*samlAssertionResponseData `json:"data"`
}

type verifyFactorRequest struct {
	AppId      string `json:"app_id"`
	DeviceId   int `json:"device_id,string"`
	StateToken string `json:"state_token"`
	OtpToken   string `json:"otp_token"`
}

type verifyFactorResponse struct {
	Status *responseStatus `json:"status"`
	Data   string `json:"data"`
}

func parseJson(resp *http.Response, v interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
