package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Username string               `bson:"username,omitempty" json:"username"`
	Email    string               `bson:"email,omitempty" json:"email"`
	Password []byte               `bson:"password,omitempty" json:"-"`
	Cart     []map[string]Product `bson:"cart,omitempty" json:"cart"`
}

type Product struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Order       int32               `bson:"order,omitempty" json:"order"`
	Name        string              `bson:"name,omitempty" json:"name"`
	Brand       string              `bson:"brand,omitempty" json:"brand"`
	Description string              `bson:"description,omitempty" json:"description"`
	Images      primitive.A         `bson:"images,omitempty" json:"images"`
	ProductType string              `bson:"product_type,omitempty" json:"product_type"`
	Options     []map[string]string `bson:"options,omitempty" json:"options"`
}

type News struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Order   int32              `bson:"order,omitempty" json:"order"`
	Title   string             `bson:"title,omitempty" json:"title"`
	Content string             `bson:"content,omitempty" json:"content"`
	Images  primitive.A        `bson:"images,omitempty" json:"images"`
	Heading string             `bson:"heading,omitempty" json:"heading"`
}

type ServiceMaster struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Order    int32              `bson:"order,omitempty" json:"order"`
	FullName string             `bson:"full_name,omitempty" json:"full_name"`
	Phone    string             `bson:"phone,omitempty" json:"phone"`
	Address  string             `bson:"address,omitempty" json:"address"`
	Role     string             `bson:"role,omitempty" json:"role"`
	Photo    string             `bson:"photo,omitempty" json:"photo"`
}

type PurchaseHistory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	KeyboardIDs  primitive.A        `bson:"keyboard_ids,omitempty" json:"keyboard_ids"`
	Quantity     int32              `bson:"quantity,omitempty" json:"quantity"`
	TotalPrice   float64            `bson:"total_price,omitempty" json:"total_price"`
	PurchaseDate string             `bson:"purchase_date,omitempty" json:"purchase_date"`
	PaymentType  string             `bson:"payment_type,omitempty" json:"payment_type"`
}

type Wiki struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Order   int32              `bson:"order,omitempty" json:"order"`
	Title   string             `bson:"title,omitempty" json:"title"`
	Content string             `bson:"content,omitempty" json:"content"`
	Images  primitive.A        `bson:"images,omitempty" json:"images"`
}
