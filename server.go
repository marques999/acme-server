package main

import (
	"fmt"
	"os"
	"log"
	"time"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/marques999/acme-server/admin"
	"github.com/marques999/acme-server/auth"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/orders"
	"github.com/marques999/acme-server/products"
)

func main() {

	envException := godotenv.Load()

	if envException != nil {
		log.Fatal(envException.Error())
	}

	psqlConnection := fmt.Sprintf("host=localhost user=%s dbname=%s sslmode=disable password=%s",
		getEnvOrDefault("POSTGRES_USER", "postgres"),
		getEnvOrDefault("POSTGRES_DB", "postgres"),
		getEnvOrDefault("POSTGRES_PASSWORD", "postgres"),
	)

	gin.SetMode(gin.ReleaseMode)
	database, connectionException := gorm.Open("postgres", psqlConnection)

	if connectionException != nil {
		log.Fatal(connectionException.Error())
	}

	defer database.Close()
	products.Migrate(database)
	orders.Migrate(database)
	customers.Migrate(database)
	middleware := getAuthenticator(database)
	router := gin.Default()
	auth.InitializeRoutes(middleware, router)
	admin.InitializeRoutes(database, router)
	customers.InitializeRoutes(database, middleware, router)
	products.InitializeRoutes(database, middleware, router)
	orders.InitializeRoutes(database, middleware, router)
	router.Run(getEnvOrDefault("ACME_HOSTNAME", ":3333"))
}

func getAuthenticator(database *gorm.DB) *jwt.GinJWTMiddleware {

	return &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Unauthorized: func(context *gin.Context, statusCode int, message string) {
			context.JSON(statusCode, gin.H{"error": message})
		},
		Authenticator: func(username string, password string, context *gin.Context) (string, bool) {
			return customers.Authenticate(database, username, password)
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

func getEnvOrDefault(variableKey string, defaultValue string) string {

	if lookupValue, exists := os.LookupEnv(variableKey); exists {
		return lookupValue
	} else {
		return defaultValue
	}
}
