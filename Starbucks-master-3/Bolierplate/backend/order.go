package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"encoding/json"
	"net/http"
	time2 "time"
)

type Order struct {
	OrderId     string     `json:"id" bson:"_id"`
	UserName    string     `json:"username" bson:"username"`
	Location    string     `json:"location" bson:"location"`
	Items       []Item     `json:"items" bson:"items"`
	Status      string     `json:"status" bson:"status"`
	Message     string     `json:"message" bson:"message"`
	Links       Links      `json:"links" bson:"links"`
	TotalAmount int    `json:"totalAmount" bson:"totalAmount"`
	OrderDate   time2.Time `json: "orderDate" bson: "orderDate"`
	PaymentDate time2.Time `json: "paymentDate" bson: "paymentDate"`
}

func (oc OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside createorder	")
	var o Order
	// Populate the order data from request.body to Order object
	json.NewDecoder(r.Body).Decode(&o)

	fmt.Println(o.Items[0].Name)
	fmt.Println(o.Items[1].Name)
	fmt.Println(o.Items[0])
	// Add an Id, using uuid for
	var orderId uuid.UUID
	orderId, _ = uuid.NewV4()
	o.OrderId = orderId.String()

	o.Status = "PLACED"
	o.Message = "Order has been placed"
	o.TotalAmount = 0

	for index := 0; index < len(o.Items); index += 1 {
		o.TotalAmount += o.Items[index].Price * o.Items[index].Quantity
	}

	o.OrderDate = time2.Now()

	// Write the user to mongo
	oc.session.DB("test").C("Order").Insert(&o)

	// Write content-type, statuscode, payload
	fmt.Println("New Order Created, Order ID:", o.OrderId)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(o)
}

func (oc OrderController) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	fmt.Println(username)
	var orders []Order
	iter := oc.session.DB("test").C("Order").Find(bson.M{"username": username}).Iter()
	result := Order{}
	for iter.Next(&result) {
		orders = append(orders, result)
	}

	for _, order := range orders {
		//fmt.Println(order.OrderId)
		//fmt.Println(order.Items[0])
		fmt.Println("--- ", order.OrderId)
		fmt.Println("---", order.Location)
		//fmt.Println("------------",order.Items)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&orders)

}
