package application

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/admin"
	"github.com/marques999/acme-server/auth"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/orders"
	"github.com/marques999/acme-server/products"
)

func Run() {

	// Global Configuration

	gin.SetMode(gin.ReleaseMode)

	// Database Connection

	database, connectionException := gorm.Open(common.ConnectionType, common.ConnectionString)

	if connectionException != nil {
		log.Fatal(connectionException.Error())
	}

	defer database.Close()

	// Database Migrations

	products.Migrate(database)
	orders.Migrate(database)
	customers.Migrate(database)

	// Initialize Middleware

	middleware := GetAuthenticator(database)

	// Initialize Routes

	router := gin.Default()
	auth.InitializeRoutes(middleware, router)
	admin.InitializeRoutes(database, router)
	customers.InitializeRoutes(database, middleware, router)
	products.InitializeRoutes(database, middleware, router)
	orders.InitializeRoutes(database, middleware, router)
	router.Run(common.Hostname)
}
