package server

import (
	"bytes"
	"context"
	"contract_ease/internal/domain"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ZitadelClient struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	ServicePAT   string
	HTTPClient   *http.Client
}

func NewZitadelClient(issuer, clientID, clientSecret, pat string) ZitadelClient {
	return ZitadelClient{
		IssuerURL:    issuer,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		ServicePAT:   pat,
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c ZitadelClient) CreateUser(ctx context.Context, params domain.CreateUserParams) (string, error) {
	body := map[string]any{
		"human": map[string]any{
			"userName": params.Username,
			"profile": map[string]string{
				"givenName":  params.FirstName,
				"familyName": params.LastName,
			},
			"email": map[string]any{
				"email":      params.Email,
				"isVerified": false,
			},
			"password": map[string]any{
				"password":       params.Password,
				"changeRequired": false,
			},
		},
	}
	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.IssuerURL+"/v2/users/new", bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.ServicePAT)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("zitadel CreateUser failed with status: %s", resp.Status)
	}

	var res struct {
		UserID       string `json:"id"`
		CreationDate string `json:"creationDate"`
		EmailCode    string `json:"emailCode"`
		PhoneCode    string `json:"phoneCode"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.UserID, nil
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (c ZitadelClient) AuthenticateROPC(ctx context.Context, email, password string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", email)
	data.Set("password", password)
	data.Set("scope", "openid email profile")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.IssuerURL+"/oauth/v2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Basic Auth: base64(client_id:client_secret)
	auth := base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(c.ClientID) + ":" + url.QueryEscape(c.ClientSecret)))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zitadel ROPC auth failed: %s", resp.Status)
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}
