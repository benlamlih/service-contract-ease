package server

import (
	"bytes"
	"context"
	"contract_ease/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
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
	slog.Info("creating user in zitadel",
		"username", params.Username,
		"email", params.Email)

	body := map[string]any{
		"username": params.Username,
		"profile": map[string]string{
			"givenName":  params.FirstName,
			"familyName": params.LastName,
		},
		"email": map[string]string{
			"email": params.Email,
		},
		"password": map[string]any{
			"password":       params.Password,
			"changeRequired": false,
		},
	}

	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.IssuerURL+"/v2/users/human", bytes.NewReader(jsonData))
	if err != nil {
		slog.Error("failed to create zitadel request",
			"error", err,
			"username", params.Username)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.ServicePAT)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		slog.Error("failed to execute zitadel request",
			"error", err,
			"username", params.Username)
		return "", err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = fmt.Errorf("error closing response body: %w", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		slog.Error("zitadel request failed with non-200 status",
			"status", resp.Status,
			"body", string(b),
			"username", params.Username)
		return "", fmt.Errorf("zitadel CreateUser failed: %s | message: %s", resp.Status, string(b))
	}

	var res struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		slog.Error("failed to decode zitadel response",
			"error", err,
			"username", params.Username)
		return "", err
	}

	slog.Debug("successfully created user in zitadel",
		"userId", res.UserID,
		"username", params.Username)

	return res.UserID, nil
}
