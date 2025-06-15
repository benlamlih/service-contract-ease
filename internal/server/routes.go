package server

import (
	"contract_ease/internal/repository"
	"log/slog"
	"net/http"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"

	"contract_ease/internal/config"
	"contract_ease/internal/domain"
)

func (s *Server) RegisterRoutes(tp trace.TracerProvider) http.Handler {

	r := gin.Default()
	r.Use(otelgin.Middleware("contract_ease", otelgin.WithTracerProvider(tp)))

	cfg := config.LoadConfig()
	frontendURL := cfg.App.FrontendURL

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/health", s.healthHandler)
		api.GET("/docs", s.apiReferenceHandler)
		api.POST("/auth/sign-up", s.signUpHandler)
	}

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health(c))
}

func (s *Server) apiReferenceHandler(c *gin.Context) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./docs/build/openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "ContractEase API",
		},
		Theme:    "deepSpace",
		Layout:   "modern",
		DarkMode: true,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API documentation"})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, htmlContent)
}

func (s *Server) signUpHandler(c *gin.Context) {
	var userRequest struct {
		Email     string `json:"email" binding:"required,email"`
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Role      string `json:"role" binding:"required,oneof=owner member"`
	}

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		slog.Error("failed to bind signup request",
			"error", err,
			"path", "/api/auth/sign-up")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userRequest data"})
		return
	}

	slog.Info("processing signup request",
		"email", userRequest.Email,
		"firstName", userRequest.FirstName,
		"lastName", userRequest.LastName,
		"role", userRequest.Role)

	username := domain.GenerateUsername(userRequest.FirstName, userRequest.LastName)
	params := domain.CreateUserParams{
		Username:  username,
		Email:     userRequest.Email,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Password:  userRequest.Password,
	}

	zitadelUserID, err := s.ZitadelClient.CreateUser(c, params)
	if err != nil {
		slog.Error("zitadel user creation failed",
			"error", err,
			"email", userRequest.Email,
			"username", username)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in authentication service"})
		return
	}

	slog.Info("zitadel user created successfully",
		"zitadelUserId", zitadelUserID,
		"email", userRequest.Email)

	err = s.store.CreateUser(c, repository.CreateUserParams{
		ZitadelID: zitadelUserID,
		FirstName: &userRequest.FirstName,
		LastName:  &userRequest.LastName,
		Username:  &username,
		Email:     userRequest.Email,
	})
	if err != nil {
		slog.Error("local user creation failed",
			"error", err,
			"email", userRequest.Email,
			"zitadelUserId", zitadelUserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in local database"})
		return
	}

	slog.Info("user signup completed successfully",
		"userId", zitadelUserID,
		"email", userRequest.Email,
		"username", username)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"userId": zitadelUserID,
	})
}
