package main

import (
	"Tugas2/structs"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "irvandy2"
	password = "koinworks"
	dbname   = "postgres"
)

func main() {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("normal")
	createOrder(db)
}
func createOrder(db *sql.DB) {
	var order = structs.Order{}
	sqlStatement := "insert Into orders(order_id,ordered_at,customer_name) Values($1,$2,$3) returning*"
	err := db.QueryRow(sqlStatement, "order_1", "today", "test1").Scan(&order.Order_id, &order.Ordered_at, &order.Customer_name)
	if err != nil {
		panic(err)
	}

	fmt.Println("New data", order)
}
