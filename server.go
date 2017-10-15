package main

import (
	"fmt"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/marques999/acme-backend/handler/customers"
	"github.com/marques999/acme-backend/handler/products"
	"github.com/marques999/acme-backend/model"
	"log"
	"net/http"
	"time"
)

func main() {
	(&App{}).Run(":3333")
}

type App struct {
	Router   *gin.Engine
	Database *gorm.DB
}

func (app *App) Run(host string) {

	gin.SetMode(gin.ReleaseMode)

	psqlConnection := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		"127.0.0.1",
		"marques999",
		"acme",
		"disable",
		"14191091")

	connection, connectionException := gorm.Open("postgres", psqlConnection)

	if connectionException != nil {
		log.Fatal(connectionException.Error())
	}

	defer connection.Close()
	app.Router = gin.Default()
	app.Database = model.Migrate(connection)

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:         "test zone",
		Key:           []byte("secret key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		Authenticator: app.AuthenticatorCallback,
		Unauthorized:  UnauthorizedCallback,
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}

	app.Router.POST("/login", authMiddleware.LoginHandler)

	auth := app.Router.Group("/auth")

	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	customerRoutes := app.Router.Group("/customers")
	{
		customerRoutes.GET("/", func(context *gin.Context) {
			customers.List(context, app.Database, fmt.Sprint(jwt.ExtractClaims(context)["id"]))
		})

		customerRoutes.POST("/", func(context *gin.Context) {
			customers.Insert(context, app.Database)
		})

		customerRoutes.GET("/:id", func(context *gin.Context) {
			customers.Find(context, app.Database, fmt.Sprint(jwt.ExtractClaims(context)["id"]))
		})

		customerRoutes.PUT("/:id", func(context *gin.Context) {
			customers.Update(context, app.Database, fmt.Sprint(jwt.ExtractClaims(context)["id"]))
		})

		customerRoutes.DELETE("/:id", func(context *gin.Context) {
			customers.Delete(context, app.Database, fmt.Sprint(jwt.ExtractClaims(context)["id"]))
		})
	}

	productRoutes := app.Router.Group("/products")
	{
		productRoutes.Use(authMiddleware.MiddlewareFunc())
		{
			productRoutes.POST("/", func(context *gin.Context) {
				if jwt.ExtractClaims(context)["id"] == "admin" {
					products.Insert(context, app.Database)
				} else {
					context.JSON(http.StatusUnauthorized, nil)
				}
			})

			productRoutes.GET("/:id", func(context *gin.Context) {
				products.Find(context, app.Database)
			})

			productRoutes.PUT("/:id", func(context *gin.Context) {
				if jwt.ExtractClaims(context)["id"] == "admin" {
					products.Update(context, app.Database)
				} else {
					context.JSON(http.StatusUnauthorized, nil)
				}
			})

			productRoutes.DELETE("/:id", func(context *gin.Context) {
				if jwt.ExtractClaims(context)["id"] == "admin" {
					products.Delete(context, app.Database)
				} else {
					context.JSON(http.StatusUnauthorized, nil)
				}
			})
		}

		productRoutes.GET("/", func(context *gin.Context) {
			products.List(context, app.Database)
		})
	}

	app.Router.Run(host)
}

func (app *App) AuthenticatorCallback(username string, password string, context *gin.Context) (string, bool) {

	customer := customers.GetCustomer(app.Database, username)

	if customer == nil {
		return username, false
	}

	if (username == "admin" && password == "admin") || (customer.Username == username && customer.Password == password) {
		return username, true
	}

	return username, false
}

func UnauthorizedCallback(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}
