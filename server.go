package main

import (
	"os"
	"fmt"
	"log"
	"time"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/appleboy/gin-jwt"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/orders"
	"github.com/marques999/acme-server/products"
	"github.com/marques999/acme-server/customers"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	if errors := godotenv.Load(); errors != nil {
		log.Fatal(errors.Error())
	}

	database := sqlx.MustConnect("postgres", fmt.Sprintf(
		"host=localhost user=%s dbname=%s sslmode=disable password=%s",
		getEnvOrDefault("POSTGRES_USER", "postgres"),
		getEnvOrDefault("POSTGRES_DB", "postgres"),
		getEnvOrDefault("POSTGRES_PASSWORD", "postgres"),
	))

	defer database.Close()
	customers.Migrate(database)
	products.Migrate(database)
	orders.Migrate(database)
	router := gin.Default()

	router.Use(common.CorsMiddleware(common.CorsConfig{
		Origins:         "*",
		ExposedHeaders:  "",
		Credentials:     true,
		ValidateHeaders: false,
		MaxAge:          50 * time.Second,
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
	}))

	middleware := &jwt.GinJWTMiddleware{
		Realm:      common.AuthenticationRealm,
		Key:        []byte(common.RamenRecipe),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Unauthorized: func(context *gin.Context, statusCode int, message string) {
			context.JSON(statusCode, gin.H{"error": message})
		},
		Authenticator: func(username string, password string, context *gin.Context) (string, bool) {
			return customers.Login(database, username, password)
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}

	common.InitializeRoutes(middleware, router)
	customers.InitializeRoutes(database, middleware, router)
	products.InitializeRoutes(database, middleware, router)
	orders.InitializeRoutes(database, middleware, router)
	router.Run(getEnvOrDefault("ACME_HOSTNAME", ":3333"))
}

func getEnvOrDefault(variableKey string, defaultValue string) string {

	if lookup, exists := os.LookupEnv(variableKey); exists {
		return lookup
	} else {
		return defaultValue
	}
}