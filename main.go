package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func createTable(db *sql.DB) error {
	// SQL statement to create the 'products' table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		category VARCHAR(255) NOT NULL
	);
	`
	_, err := db.Exec(createTableSQL)
	return err
}

func createProduct(db *sql.DB) {

}

func updateProduct(db *sql.DB) {

}

func getAllProducts(db *sql.DB) {

}

func deleteProduct(db *sql.DB) {

}

func getProductById(db *sql.DB) {

}

func getProductByCategory(db *sql.DB) {

}

func main() {
	user := "root"
	password := "Manager0"
	host := "127.0.0.1"
	port := 3306
	dbName := "mysql"
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mysql db...")

	err = createTable(db)
	if err != nil {
		log.Fatal("Error creating 'products' table:", err)
	}

	fmt.Println("Table 'products' created successfully!")

	var (
		id   int
		name string
	)

	//Insert a new product

	res, err := db.Exec("insert into products (name) values (?)", "Macbook m2")
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ID=%d, affected =%d", lastId, rowCnt)

	//Query the db

	id = 1

	rows, err := db.Query("select id, name from products where id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
