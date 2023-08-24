package database

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/MadiSergazy/mongodb/GolangEcommerceProject/models"
)

var (
	ErrCantFindProduct    = errors.New("can't find product")
	ErrCantDecodeProducts = errors.New("can't find product")
	ErrUserIDIsNotValid   = errors.New("user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add product to cart")
	ErrCantRemoveItem     = errors.New("cannot remove item from cart")
	ErrCantGetItem        = errors.New("cannot get item from cart ")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDB, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	// ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	// defer cancel()
	var productCarts []models.ProductUser
	for searchFromDB.Next(ctx) {
		var productCart models.ProductUser
		if err = searchFromDB.Decode(&productCart); err != nil {
			log.Println(err)
			return ErrCantDecodeProducts
		}

		productCarts = append(productCarts, productCart)
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCarts}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	return nil
}

func RemoveCartItem(ctx context.Context, prodCOllection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)
	// db.userCollection.updateMany(
	// 	{ _id: id },
	// 	{ $pull: { usercart: { _id: productID } } }
	//  );
	if err != nil {
		return ErrCantRemoveItem
	}

	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	// 1 fetch cart from user
	// 2 find the cart total
	// 3 create order with the items
	// 4 added order to the use collection
	// 5 added cart item in the cart to order list
	// 6 empty up the cart

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}
	var getCartItems models.User
	var orderCart models.Order

	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Orderered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	// * Step 1: Unwind the usercart array
	//path is going to tell what exactly you want to unwind
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	//using unwind we can get accsess to all fields

	// ^ Step 2: Group to calculate the total price of the cart items
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$id"}, {Key: "total", Value: bson.D{primitive.E{Key: "sum", Value: "$usercart.price"}}}}}}
	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})

	// db.userCollection.aggregate([
	// 	{ $unwind: "$usercart" },
	// 	{
	// 	  $group: {
	// 		_id: "$_id",
	// 		total: { $sum: "$usercart.price" }
	// 	  }
	// 	}
	//   ]);

	if err != nil {
		log.Println(err)
		return err
	}
	defer ctx.Done()

	var getUserCarts bson.M
	if err = currentResults.All(ctx, &getUserCarts); err != nil {
		log.Println(err)
		return err
	}

	var total_price int32
	for _, user_item := range getUserCarts {
		if userMap, ok := user_item.(map[string]interface{}); ok {
			if price, exists := userMap["total"]; exists {
				total_price = price.(int32)
			}
		}
	}
	// ^Step 3: Add order to the user collection
	orderCart.Price = int(total_price)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}

	// db.userCollection.updateMany(
	// 	{ _id: id },
	// 	{ $push: { orders: orderCart } }
	//  );

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return err
	}
	if err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems); err != nil {
		log.Println(err)
		return err
	}
	//^ { $push: { "orders.$[].order_list": { $each: getCartItems.UserCart } } }: This update operation uses the $push operator to add the items from getCartItems.
	//^ UserCart into the order_list array of all orders within the orders array. The $[] operator is used to update all elements in the orders array.
	// * Step 5: Add cart items to the orders_list of the order
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}

	// db.userCollection.updateOne(
	// 	{ _id: id },
	// 	{ $push: { "orders.$[].order_list": { $each: getCartItems.UserCart } } }
	//  );

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
		return err
	}

	// * 4 empty up the cart
	// *Step 6: Empty the usercart
	usercart_empty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}

	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}

	var productDetails models.ProductUser
	var orderDetail models.Order

	orderDetail.Order_ID = primitive.NewObjectID()
	orderDetail.Orderered_At = time.Now()
	orderDetail.Order_Cart = make([]models.ProductUser, 0)
	orderDetail.Payment_Method.COD = true

	if err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails); err != nil {
		log.Println(err)
		return err
	}
	orderDetail.Price = productDetails.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderDetail}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return err
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
