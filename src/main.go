package main

import (
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

var DB *gorm.DB

func initDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&Product{})
}

func main() {
	app := fiber.New()

	initDatabase()

	app.Post("/products", createProduct)
	app.Get("/products", getProducts)
	app.Get("/products/:id", getProduct)
	app.Put("/products/:id", updateProduct)
	app.Delete("/products/:id", deleteProduct)

	app.Listen(":3000")
}

func createProduct(c *fiber.Ctx) error {
	product := new(Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	DB.Create(product)
	return c.JSON(product)
}

func getProducts(c *fiber.Ctx) error {
	var products []Product
	DB.Find(&products)
	return c.JSON(products)
}

func getProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product Product
	result := DB.First(&product, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}
	return c.JSON(product)
}

func updateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product Product
	if err := DB.First(&product, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	DB.Save(&product)
	return c.JSON(product)
}

func deleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := DB.Delete(&Product{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
