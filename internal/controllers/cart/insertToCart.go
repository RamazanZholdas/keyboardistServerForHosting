package cart

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/RamazanZholdas/KeyboardistSV2/internal/app"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/jwt"
	"github.com/RamazanZholdas/KeyboardistSV2/internal/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertToCart(c *fiber.Ctx) error {
	var optionId string

	var requestBody struct {
		OptionId string `json:"optionId"`
	}

	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}
	optionId = requestBody.OptionId

	cookie := c.Cookies("jwt")

	claims, err := jwt.ExtractTokenClaimsFromCookie(cookie)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	var user models.User
	err = app.GetMongoInstance().FindOne("users", bson.M{"email": claims.Issuer}, &user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	order := c.Params("order")
	number, err := strconv.Atoi(order)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "The order must be a valid number.",
		})
	}

	fmt.Println("entering the loop")
	//check if product exists in users cart
	for _, cartItem := range user.Cart {
		stringOrderCart := strconv.Itoa(int(cartItem["product"].Order))
		if order == stringOrderCart {
			quantity, err := strconv.Atoi(cartItem["product"].Options[0]["quantity"])
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to updatSe cart"})
			}
			quantity++
			cartItem["product"].Options[0]["quantity"] = strconv.Itoa(quantity)

			update := bson.M{"$set": bson.M{"cart": user.Cart}}
			err = app.GetMongoInstance().UpdateOne("users", bson.M{"_id": user.ID}, update)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update cart"})
			}
			return c.JSON(fiber.Map{"message": "Product added to cart"})
		}
	}

	var product models.Product
	err = app.GetMongoInstance().FindOne("products", bson.M{"order": number}, &product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	filteredOptions := make([]map[string]string, 0)
	for _, option := range product.Options {
		if option["optionId"] == optionId {
			filteredOptions = append(filteredOptions, option)
		}
	}

	product.Options = filteredOptions

	// setting up quantity
	product.Options[0]["quantity"] = "1"

	user.Cart = append(user.Cart, map[string]models.Product{"product": product})

	update := bson.M{"$set": bson.M{"cart": user.Cart}}
	err = app.GetMongoInstance().UpdateOne("users", bson.M{"_id": user.ID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update cart"})
	}

	return c.JSON(fiber.Map{"message": "Product added to cart"})
}
