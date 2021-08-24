package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const route = "/api/{name}/{time}"

type Entry struct {
	ID primitive.ObjectID `bson:"_id"`
	UserTime int	`bson:"time"`
	UserDate time.Time	`bson:"date"`
	UserName string		`bson:"name"`
}

type Entries []Entry

func main() {
	var port = flag.Int("p", 8080, "Port to run the server")
	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/api", GetHandlerFunc).Methods("GET")
	r.HandleFunc(route, PostHandlerFunc).Methods("POST")

	fmt.Println("Listening on port: ", *port)
    http.ListenAndServe(fmt.Sprintf(":%v", *port), r)
}

func dbConn(ctx context.Context) (*mongo.Client, *mongo.Collection) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("Database connection error")
	}

	recordCollection := client.Database("Record").Collection("Record")

	return client, recordCollection
}

func GetHandlerFunc(w http.ResponseWriter, r *http.Request){

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, recordCollection := dbConn(ctx)

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Error: ",err)
		}
	}()

	cur, err := recordCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal("Search error: ", err)
	}
	defer cur.Close(ctx)

	
	var results = Entries{}

	for cur.Next(context.Background()){
		var res Entry

		err := cur.Decode(&res)
		if err != nil {
			log.Fatal("Decode error: ",err)
		}
		
		results = append(results, res)
	}

	if err != nil {
		log.Fatal("Json Error")
	}

	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}

func PostHandlerFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["name"]
	userTime := vars["time"]
	intTime, err := strconv.Atoi(userTime)
	
	if err != nil {
		log.Fatal("conversion error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, recordCollection := dbConn(ctx)

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Error: ",err)
		}
	}()

	_, err = recordCollection.InsertOne(ctx, bson.D{
		{"name", userName},
    	{"time", intTime},
		{"date", time.Now()},
		})
	if err != nil {
		log.Fatal("Insert error: ", err)
	}

	var resp = struct{
		status string
	}{status: "OK"}
	json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)

}