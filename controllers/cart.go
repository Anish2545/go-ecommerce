package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Anish2545/go-ecommerce/models"
	"github.com/Anish2545/go-ecommerce/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Application struct{
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection,userCollection *mongo.Collection) *Application{
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart()  gin.HandlerFunc{
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("Product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("User id is empty"))
			return
		}

		productID,err := primitive.ObjectIDFromHex(productQueryID)

		if err!=nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx,cancel = context .WithTimeout(context.Background(),5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx,app.prodCollection,app.userCollection,productID,userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,err)
		}
		c.IndentedJSON(200,"Successfully added to the cart")
	}
}

func (app *Application) RemoveItem()  gin.HandlerFunc{			
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("Product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("User id is empty"))
			return
		}

		productID,err := primitive.ObjectIDFromHex(productQueryID)

		if err!=nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx,cancel = context .WithTimeout(context.Background(),5*time.Second)
		defer cancel()


		database.RemoveCartItem(ctx,app.prodCollection,app.userCollection,productID,userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,err)
			return
		}
		c.IndentedJSON(200,"Successfuly removed Cart Item")
	}
}

func GetItemFromCart()  gin.HandlerFunc{
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == ""{
			c.Header("Content_type","application/json")
			c.JSON(http.StatusNotFound,gin.H{"Error":"invalid id"})
			c.Abort()
			return
		}

		usert_id ,_ := primitive.ObjectIDFromHex(user_id)
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var filledcart models.User
		err := userCollection.FindOne(ctx,bson.D{primitive.E{Key: "_id",Value: usert_id}}).Decode(&filledcart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500,"NOT FOUND")
			return
		}
		filter_match :=  bson.D{{Key: "$match",Value: bson.D{primitive.E{Key: "_id",Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind",Value: bson.D{primitive.E{Key: "path",Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group",Value: bson.D{primitive.E{Key: "_id",Value: "$_id"},{Key: "total",Value: bson.D{primitive.E{Key: "$sum",Value: "$usercart.price"}}}}}}

		pointcursor,err := userCollection.Aggregate(ctx,mongo.Pipeline{filter_match,unwind,grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err =pointcursor.All(ctx,&listing);err!=nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _,json := range listing{
			c.IndentedJSON(200,json["total"])
			c.IndentedJSON(200,filledcart.UserCart)
		}
		ctx.Done()

	}
}

func (app *Application) BuyFromCart()  gin.HandlerFunc{
	return func(c *gin.Context) {
		userQueryID := c.Query("id")

		if userQueryID == ""{
			log.Panic("user id is empty")
			_ =c.AbortWithError(http.StatusBadRequest,errors.New("user ID is Empty"))
		}

		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		defer cancel()

		err := database.BuyItemFromCart(ctx,app.userCollection,userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,err)
		}
		c.IndentedJSON(200,"Seccessfuly placed the order")
	}
}

func (app *Application) InstantBuy()  gin.HandlerFunc{
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("Product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest,errors.New("User id is empty"))
			return
		}

		productID,err := primitive.ObjectIDFromHex(productQueryID)

		if err!=nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx,cancel = context .WithTimeout(context.Background(),5*time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx,app.prodCollection,app.userCollection,productID,userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError,err)
		}
		c.IndentedJSON(200,"Successfuly placed the order")
	}
}


