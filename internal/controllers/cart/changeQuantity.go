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

func ChangeQuantity(c *fiber.Ctx) error {
	var requestBody struct {
		Quantity string `json:"quantity"`
	}

	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

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

	var product models.Product
	err = app.GetMongoInstance().FindOne("products", bson.M{"order": number}, &product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Product not found"})
	}

	var productFromCart models.Product
	for index, item := range user.Cart {
		productFromCart = item["product"]

		if productFromCart.Order == int32(number) {
			if productFromCart.Options[0]["inStock"] < requestBody.Quantity {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "The quantity must be less than or equal to the number of items in stock.",
				})
			}

			oldQuantity, err := strconv.Atoi(productFromCart.Options[0]["quantity"])
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
			}
			newQuantity, err := strconv.Atoi(requestBody.Quantity)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body is invalid"})
			}

			resultQuantity := oldQuantity + newQuantity
			productFromCart.Options[0]["quantity"] = strconv.Itoa(resultQuantity)

			user.Cart[index]["product"] = productFromCart

			fmt.Println(productFromCart)

			err = app.GetMongoInstance().UpdateOne("users", bson.M{"email": claims.Issuer}, bson.M{"$set": bson.M{"cart": user.Cart}})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
			}
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "The quantity of the product has been changed successfully.",
	})
}
