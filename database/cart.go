package controllers

import "errors"

var (
	ErrCantFindProduct = errors.New("Cant find the Product")
	ErrCantDecodeProducts = errors.New("Cant find the Product") 
	ErrUserIdIsNotValid = errors.New("This user is not valid")
	ErrCantUpdateUser = errors.New("Cannot Add this product to the cart")
	ErrCantRemoveItemCart = errors.New("Cannot Remove this product to the cart")
	ErrCantGetItem = errors.New("Unable to get this item from the cart")
	ErrCantBuyCartItem = errors.New("Cannot update the Purchase")
)

func AddProductToCart()  {
	
}

func RemoveCartItem()  {
	
}

func BuyItemFromCart()  {
	
}

func InstantBuyer()  {
	
}