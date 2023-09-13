package models

// CartItem representa un artículo en el carrito de compras

type CartItem struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}