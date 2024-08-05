package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Anish2545/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct = errors.New("Cant find the Product")
	ErrCantDecodeProducts = errors.New("Cant find the Product") 
	ErrUserIdIsNotValid = errors.New("This user is not valid")
	ErrCantUpdateUser = errors.New("Cannot Add this product to the cart")
	ErrCantRemoveItemCart = errors.New("Cannot Remove this product to the cart")
	ErrCantGetItem = errors.New("Unable to get this item from the cart")
	ErrCantBuyCartItem = errors.New("Cannot update the Purchase")
)

func AddProductToCart(ctx context.Context,prodCollection,userCollection *mongo.Collection, productID primitive.ObjectID,userID string) error  {
	searchfromdb,err := prodCollection.Find(ctx,bson.M{"_id":productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}	
	var productCart []models.ProductUser
	err = searchfromdb.All(ctx,&productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}
	id,err := primitive.ObjectIDFromHex(userID)
	if err!=nil{
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key:"_id",Value :id}}
	update := bson.D{{Key:"$push", Value:bson.D{primitive.E{Key:"usercart", Value:bson.D{Key:"$each",Value:productCart}}}}}
	_,err = userCollection.UpdateOne(ctx,filter,update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection,userCollection *mongo.Collection,productID primitive.ObjectID,userID string) error {
	id,err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	filter := bson.D{primitive.E{Key:"_id",Value:id}}
	update := bson.M{"$pull":bson.M{"usercart":bson.M{"_id":productID}}}
	_,err = userCollection.UpdateMany(ctx,filter,update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil
}



func BuyItemFromCart(ctx context.Context,userCollection *mongo.Collection, productID primitive.ObjectID,userID string) error {
	id,err:=primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getcartitems models.User
	var ordercart models.Order

	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD=true

	unwind := bson.D{{Key: "$unwind" , Value: bson.D{primitive.E{Key: "path",Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group" , Value: bson.D{primitive.E{Key:"_id" , Value: "$_id"},{Key: "total",Value: bson.D{primitive.E{Key: "$sum",Value: "$usercart.price"}}}}}}
	currentresult , err := userCollection.Aggregate(ctx,mongo.Pipeline{unwind,grouping})

	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getusercart []bson.M
	if err = currentresult.All(ctx,&getusercart);err!=nil{
		panic(err)
	}

	var total_price int32

	for _,user_item := range getusercart{
		price := user_item["total"]
		total_price = price.(int32)
	}

	ordercart.Price  = int(total_price)
	filter := bson.D{primitive.E{Key:"_id",Value:id}}
	update := bson.D{{Key: "$push",Value: bson.D{primitive.E{Key: "orders",Value: ordercart}}}}
	_,err = userCollection.UpdateMany(ctx,filter,update)
	if err != nil {
		log.Println(err)
	}
	
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id",Value: id}}).Decode(&getcartitems)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key:"_id",Value:id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each":getcartitems.UserCart}}}

	_,err = userCollection.UpdateOne(ctx,filter2,update2)
	if err != nil {
		log.Println(err)
	}

	usercart_empty := make([]models.ProductUser,0)
	filter3 := bson.D{primitive.E{Key:"_id",Value:id}}
	update3 := bson.D{{Key: "$set",Value: bson.D{primitive.E{Key:"usercart",Value: usercart_empty}}}}

	_,err = userCollection.UpdateOne(ctx,filter3,update3)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer()  {
	
}