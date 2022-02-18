package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type orders struct {
	Order_id      string
	Ordered_at    string
	Customer_name string
	Items         item
}
type item struct {
	Item_id     string
	Item_name   string
	Description string
	Quantity    string
	Order_id    string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "irvandy2"
	password = "koinworks"
	dbname   = "postgres"
	PORT     = ":8080"
)

var connString string
var db *sql.DB
var err error

func main() {
	connString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("normal")
	http.HandleFunc("/getorder", getOrder)
	http.HandleFunc("/createorder", createOrder)
	http.HandleFunc("/updateorder", updateOrder)
	http.HandleFunc("/deleteorder", deleteOrder)
	http.ListenAndServe(PORT, nil)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//var order = orders{}
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
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(result)
		return
	}
	http.Error(w, "Invalid Method", http.StatusBadRequest)
}
func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var order = orders{}
	if r.Method == "POST" {
		order_id := r.FormValue("order_id")
		ordered_at := r.FormValue("ordered_at")
		customer_name := r.FormValue("customer_name")
		item_id := r.FormValue("item_id")
		item_name := r.FormValue("item_name")
		item_description := r.FormValue("item_description")
		item_quantity := r.FormValue("item_quantity")
		order = orders{
			Order_id:      order_id,
			Ordered_at:    ordered_at,
			Customer_name: customer_name,
			Items: item{
				Item_id:     item_id,
				Item_name:   item_name,
				Description: item_description,
				Quantity:    item_quantity,
			},
		}
	}
	sqlStatement := "insert Into orders(order_id,ordered_at,customer_name) Values($1,$2,$3) returning*"
	err := db.QueryRow(sqlStatement, order.Order_id, order.Ordered_at, order.Customer_name).Scan(&order.Order_id, &order.Ordered_at, &order.Customer_name)
	if err != nil {
		panic(err)
	}
	sqlStatement2 := "insert Into items(item_id,item_code,description,quantity,order_id) Values($1,$2,$3,$4,$5) returning*"
	err2 := db.QueryRow(sqlStatement2, order.Items.Item_id, order.Items.Item_name, order.Items.Description, order.Items.Quantity, order.Order_id).Scan(&order.Items.Item_id, &order.Items.Item_name, &order.Items.Description, &order.Items.Quantity, &order.Order_id)
	if err2 != nil {
		panic(err)
	}
	fmt.Println("New data", order)
}
func deleteOrder(w http.ResponseWriter, r *http.Request) {
	var order_id string
	if r.Method == "DELETE" {
		order_id = r.FormValue("order_id")
	}
	sqlStatement := "delete from orders where order_id = $1"
	res, err := db.Exec(sqlStatement, order_id)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("deleted", count)
}
func updateOrder(w http.ResponseWriter, r *http.Request) {
	var order = orders{}
	sqlStatement := "update orders set ordered_at = $2, customer_name = $3 where order_id = $1"
	if r.Method == "PUT" {
		order_id := r.FormValue("order_id")
		ordered_at := r.FormValue("ordered_at")
		customer_name := r.FormValue("customer_name")

		order = orders{
			Order_id:      order_id,
			Ordered_at:    ordered_at,
			Customer_name: customer_name,
		}
		res, err := db.Exec(sqlStatement, order.Order_id, ordered_at, customer_name)
		if err != nil {
			panic(err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(order)
		fmt.Println("affected", count)
		return
	}

}
