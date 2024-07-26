package routes

import(
	"github.com/Anish2545/go-ecommerce/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/users/signup",controllers.Signup())
	incomingRoutes.POST("/users/login",controllers.Login())
	incomingRoutes.POST("/admin/addproduct",controllers.ProductViewAdmin())
	incomingRoutes.GET("/users/productview",controllers.SearchProduct())
	incomingRoutes.GET("/users/search",controllers.SearchProductByQuery())
}