package main

import (
	"fmt"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/products"
	"github.com/marques999/acme-server/tokens"
	_ "github.com/mattn/go-sqlite3"
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

func (app *App) Migrate(database *gorm.DB) *gorm.DB {
	products.Migrate(database)
	tokens.Migrate(database)
	customers.Migrate(database)
	return database
}

func (app *App) Run(host string) {

	gin.SetMode(gin.ReleaseMode)

	/*psqlConnection := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		"127.0.0.1",
		"marques999",
		"acme",
		"disable",
		"14191091")

	connection, connectionException := gorm.Open("postgres", psqlConnection)*/
	connection, connectionException := gorm.Open("sqlite3", "acme.db")

	if connectionException != nil {
		log.Fatal(connectionException.Error())
	}

	defer connection.Close()
	app.Router = gin.Default()
	app.Database = app.Migrate(connection)

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(username string, password string, context *gin.Context) (string, bool) {
			return customers.Authenticate(app.Database, username, password)
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
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
		customerRoutes.Use(authMiddleware.MiddlewareFunc())
		{
			customerRoutes.GET("/", func(context *gin.Context) {
				if jwt.ExtractClaims(context)["id"] == "admin" {
					customers.List(context, app.Database)
				} else {
					context.JSON(http.StatusUnauthorized, nil)
				}
			})

			customerRoutes.GET("/:id", func(context *gin.Context) {
				customers.Find(context, app.Database, fmt.Sprint(jwt.ExtractClaims(context)["id"]))
			})

			customerRoutes.POST("/", func(context *gin.Context) {
				customers.Insert(context, app.Database)
			})

			customerRoutes.PUT("/:id", func(context *gin.Context) {
				customers.Update(context, app.Database, fmt.Sprint(jwt.ExtractClaims(context)["id"]))
			})

			customerRoutes.DELETE("/:id", func(context *gin.Context) {
				customers.Delete(context, app.Database, fmt.Sprint(jwt.ExtractClaims(context)["id"]))
			})
		}
	}

	productRoutes := app.Router.Group("/products")
	{
		productRoutes.Use(authMiddleware.MiddlewareFunc())
		{
			productRoutes.GET("/", func(context *gin.Context) {
				products.List(context, app.Database)
			})

			productRoutes.GET("/:id", func(context *gin.Context) {
				products.Find(context, app.Database)
			})

			productRoutes.POST("/", func(context *gin.Context) {

				if jwt.ExtractClaims(context)["id"] == "admin" {
					products.Insert(context, app.Database)
				} else {
					context.JSON(http.StatusUnauthorized, nil)
				}
			})
		}
	}

	tokenRoutes := app.Router.Group("/tokens")
	{
		tokenRoutes.Use(authMiddleware.MiddlewareFunc())
		{
			tokenRoutes.GET("/", func(context *gin.Context) {
				tokens.List(context, app.Database, fmt.Sprint(jwt.ExtractClaims(context)["id"]))
			})

			tokenRoutes.GET("/:id", func(context *gin.Context) {
				tokens.Find(context, app.Database, fmt.Sprint(jwt.ExtractClaims(context)["id"]))
			})

			tokenRoutes.POST("/", func(context *gin.Context) {
				tokens.Insert(context, app.Database)
			})

			tokenRoutes.DELETE("/:id", func(context *gin.Context) {
				tokens.Delete(context, app.Database, fmt.Sprint(jwt.ExtractClaims(context)["id"]))
			})
		}
	}

	app.Router.Run(host)
}
