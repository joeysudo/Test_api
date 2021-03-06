package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Order struct (Model)
type Order struct {
	OrderID         string    `json:"order_id"`
	OrderName       string    `json:"order_name"`
	CreateTime      time.Time `json:"create_time"`
	DeliveredAmount float64   `json:"delivered_amount"`
	TotalAmount     float64   `json:"total_amount"`
	Name        string `json:"customer_name"`
	CompanyName string `json:"company_name"`
}



func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome testorder!")
}

type server struct{}

// Init orders var as a slice Order struct
var orders []Order

// Get all orders
func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(orders)
}

// Get single order by orderID
func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r) // Gets params
	// Loop through orders and find one with the id from the params
	for _, item := range orders {
		if item.OrderID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Order{})
}

// Add new order
func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var order Order
	_ = json.NewDecoder(r.Body).Decode(&order)
	orders = append(orders, order)
	json.NewEncoder(w).Encode(order)
}

// Update order
func updateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	for index, item := range orders {
		if item.OrderID == params["id"] {
			orders = append(orders[:index], orders[index+1:]...)
			var order Order
			_ = json.NewDecoder(r.Body).Decode(&order)
			order.OrderID = params["id"]
			orders = append(orders, order)
			json.NewEncoder(w).Encode(order)
			return
		}
	}
	json.NewEncoder(w).Encode(orders)
}

// Delete order
func deleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	for index, item := range orders {
		if item.OrderID == params["id"] {
			orders = append(orders[:index], orders[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(orders)
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8000"
}

// Main function
func main() {

	// @todo: add database
	csvFile, _ := os.Open("test_data.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		record, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		da, _ := strconv.ParseFloat(record[15], 64)
		ta, _ := strconv.ParseFloat(record[16], 64)
		t, _ := time.Parse(time.RFC3339 , record[1])
		orders = append(orders, Order{
			OrderID:         record[0],
			OrderName:       record[2]+record[7],
			CreateTime:      t,
			DeliveredAmount: da,
			TotalAmount:     ta,
			Name:        record[11],
			CompanyName: record[13],
		})
	}
	orderJSON, _ := json.Marshal(orders)
	fmt.Println(string(orderJSON))
	fmt.Println("Listen to localhost:8000!")
	// Init router
	r := mux.NewRouter().StrictSlash(true)
	// Route handles & endpoints
	r.HandleFunc("/", homeLink).Methods("GET")
	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", updateOrder).Methods("PUT")
	r.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")
	// Start server
	port := getPort()
	log.Fatal(http.ListenAndServe(port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))

}
