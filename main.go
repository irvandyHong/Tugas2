package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type orders struct {
	Order_id      string
	Ordered_at    string
	Customer_name string
}

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
	//getOrder(db)
	//deleteOrder(db)
	updateOrder(db)
	getOrder(db)

}
func createOrder(db *sql.DB) {
	var order = orders{}
	sqlStatement := "insert Into orders(order_id,ordered_at,customer_name) Values($1,$2,$3) returning*"
	err := db.QueryRow(sqlStatement, "order_1", "today", "test1").Scan(&order.Order_id, &order.Ordered_at, &order.Customer_name)
	if err != nil {
		panic(err)
	}

	fmt.Println("New data", order)
}
func getOrder(db *sql.DB) {
	var result = []orders{}
	sqlStatement := "select * from orders"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var order = orders{}

		err = rows.Scan(&order.Order_id, &order.Ordered_at, &order.Customer_name)
		if err != nil {
			panic(err)
		}
		result = append(result, order)
	}
	fmt.Println("Order Data", result)
}
func deleteOrder(db *sql.DB) {
	sqlStatement := "delete from orders where order_id = $1"
	res, err := db.Exec(sqlStatement, "order_1")
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("deleted", count)
}
func updateOrder(db *sql.DB) {
	sqlStatement := "update orders set ordered_at = $2, customer_name = $3 where order_id = $1"
	res, err := db.Exec(sqlStatement, "order_1", "tomorrow", "cust2")
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("affected", count)
}
