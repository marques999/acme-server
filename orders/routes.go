package orders

import (
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
	"github.com/marques999/acme-server/common"
)

func InitializeRoutes(database *sqlx.DB, middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	routes := router.Group("/" + Orders)
	{
		routes.Use(middleware.MiddlewareFunc())
		{
			routes.GET(common.RouteDefault, func(context *gin.Context) {
				context.JSON(List(database, common.ParseId(context)))
			})

			routes.POST(common.RouteDefault, func(context *gin.Context) {
				context.JSON(Insert(context, database, common.ParseId(context)))
			})

			routes.GET(common.RouteWithId, func(context *gin.Context) {
				context.JSON(Find(context, database, common.ParseId(context)))
			})

			routes.PUT(common.RouteWithId, func(context *gin.Context) {
				context.JSON(Purchase(context, database, common.ParseId(context)))
			})

			routes.DELETE(common.RouteWithId, func(context *gin.Context) {
				context.JSON(Delete(context, database, common.ParseId(context)))
			})
		}
	}
}