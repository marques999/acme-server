package customers

import (
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
	"github.com/marques999/acme-server/common"
)

func InitializeRoutes(database *sqlx.DB, middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	routes := router.Group("/" + Customers)
	{
		routes.POST(common.RouteDefault, func(context *gin.Context) {
			context.JSON(Post(context, database))
		})

		routes.Use(middleware.MiddlewareFunc())
		{
			routes.GET(common.RouteDefault, func(context *gin.Context) {
				context.JSON(List(database, common.ParseId(context)))
			})

			routes.GET(common.RouteWithId, func(context *gin.Context) {
				context.JSON(Find(context, database, common.ParseId(context)))
			})

			routes.PUT(common.RouteWithId, func(context *gin.Context) {
				context.JSON(Put(context, database, common.ParseId(context)))
			})

			routes.DELETE(common.RouteWithId, func(context *gin.Context) {
				context.JSON(Delete(context, database, common.ParseId(context)))
			})
		}
	}
}