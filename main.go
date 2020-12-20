	package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"strconv"
	//"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

type User struct {
	Id int `json:"Id" bson:"Id"`
	Name string `json:"Name" bson:"Name"`
	Date_of_birth string `json:"Date_of_birth" bson:"Date_of_birth"`
	Phone_number string `json:"Phone_number" bson:"Phone_number"`
	Email_address string `json:"Email_address" bson:"Email_address"`
	Created time.Time `json:"Created" bson:"Created"`
}

type Contact struct {
	UserIdOne int `json:"UserIdOne"`
	UserIdTwo int `json:"UserIdTwo"`
	Timestamp time.Time `json:"Timestamp"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		userid, _ := strconv.Atoi(r.URL.Query().Get("Id"))
		fmt.Print(userid)
		filter := bson.D{{"Id",userid}}
		collection := client.Database("appointytask").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		var user User
		result := collection.FindOne(ctx, filter).Decode(&user)
		fmt.Print(result)
		jsonData, _ := json.Marshal(result)
		w.Write(jsonData)
	}else {
		w.Header().Set("Content-Type", "application/json")
		reqBody, _ := ioutil.ReadAll(r.Body)
	    var user User
	    json.Unmarshal(reqBody, &user)
		collection := client.Database("appointytask").Collection("users")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_ , _ = collection.InsertOne(ctx, user)
		jsonData, _ := json.Marshal(user)
		w.Write(jsonData)
	}
}

func createContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, _ := ioutil.ReadAll(r.Body)
    var contact Contact
    json.Unmarshal(reqBody, &contact)
	collection := client.Database("appointytask").Collection("contacts")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_ , _ = collection.InsertOne(ctx, contact)
	jsonData, _ := json.Marshal(contact)
	w.Write(jsonData)
}

func findUser(w http.ResponseWriter, r *http.Request) {
	
}

func handleRequests() {
	http.HandleFunc("/users", createUser)
	http.HandleFunc("/contacts", createContact)
	http.HandleFunc("/users/:Id", findUser)
	log.Fatal(http.ListenAndServe(":5000",nil))
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	handleRequests()
}