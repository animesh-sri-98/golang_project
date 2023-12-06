package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	ID         int
	Name       string
	Price      float64
	CategoryID int
}
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Get all products
func getProducts(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CategoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		products = append(products, p)

	}
	c.JSON(http.StatusOK, products)
}

// Get a single product by ID
func getProductsById(c *gin.Context) {
	str_id := c.Param("id")
	id, err := strconv.Atoi(str_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	var p Product
	result := db.QueryRow("SELECT * from products WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Price, &p.CategoryID)

	if result == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	} else if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)

}

func getProductByCategory(c *gin.Context) {
	category_str := c.Param("category_id")
	category_id, err := strconv.Atoi(category_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, err := db.Query("SELECT * FROM products WHERE category_id = ?", category_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var products []Product

	for rows.Next() {
		var p Product

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CategoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
		c.JSON(http.StatusOK, products)

	}

}

// Create a new product
func createProduct(c *gin.Context) {
	var product Product

	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO products (name, price, category_id) VALUES (?, ?, ?)", product.Name, product.Price, product.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	productID, _ := result.LastInsertId()
	product.ID = int(productID)

	c.JSON(http.StatusCreated, product)
}

// Update a product by ID
func updateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var updatedProduct Product

	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec("UPDATE products SET name = ?, price = ?, category_id = ? WHERE id = ?", updatedProduct.Name, updatedProduct.Price, updatedProduct.CategoryID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func deleteProduct(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Delete the product from the database
	_, err = db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// Get all categories
func getCategories(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var cat Category
		err := rows.Scan(&cat.ID, &cat.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		categories = append(categories, cat)
	}
	c.JSON(http.StatusOK, categories)

}

// Get a single category by ID
func getCategoryById(c *gin.Context) {
	// Get the category ID from the request parameters
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the category from the database by ID
	var category Category
	err = db.QueryRow("SELECT * FROM categories WHERE id = ?", id).Scan(&category.ID, &category.Name)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Create a new category
func createCategory(c *gin.Context) {
	// Parse the JSON request body to get category data
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the new category into the database
	result, err := db.Exec("INSERT INTO categories (name) VALUES (?)", category.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categoryID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the ID in the category struct
	category.ID = int(categoryID)

	// Return the created category as JSON
	c.JSON(http.StatusCreated, category)
}

// Update a category by ID
func updateCategory(c *gin.Context) {
	// Get the category ID from the request parameters
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Parse the JSON request body to get updated category data
	var updatedCategory Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the category in the database
	_, err = db.Exec("UPDATE categories SET name = ? WHERE id = ?", updatedCategory.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success status
	c.Status(http.StatusNoContent)
}

// Delete a category by ID
func deleteCategory(c *gin.Context) {
	// Get the category ID from the request parameters
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Delete the category from the database by ID
	_, err = db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success status
	c.Status(http.StatusNoContent)
}

func initDB(dataSourceName string) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mysql db...")
	createTables(db)

}

func createTables(db *sql.DB) {
	// Create categories table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create products table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			category_id INT,
			FOREIGN KEY (category_id) REFERENCES categories(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Tables created successfully")
}

var db *sql.DB

func main() {

	user := "root"
	password := "Manager0"
	host := "127.0.0.1"
	port := 3306
	dbName := "mysql"
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName)
	initDB(dataSourceName)

	//router
	router := gin.Default()

	router.GET("/products", getProducts)
	router.GET("/products/:id", getProductsById)
	router.POST("/products", createProduct)
	router.PUT("/products/:id", updateProduct)
	router.DELETE("/products/:id", deleteProduct)
	router.GET("/products/category", getProductByCategory)

	router.GET("/categories", getCategories)
	router.GET("/categories/:id", getCategoryById)
	router.POST("/categories", createCategory)
	router.PUT("/categories/:id", updateCategory)
	router.DELETE("/categories/:id", deleteCategory)

	server_port := 8080
	fmt.Println("Server started on port..", server_port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server_port), router))
}
