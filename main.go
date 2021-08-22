package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const route = "/{name}/{date}/{time}"

type records struct {
	client 	*mongo.Client
	ctx 	context.Context
}

type Entry struct {
	ID primitive.ObjectID `bson:"_id"`
	UserTime int	`bson:"time"`
	UserDate time.Time	`bson:"date"`
}

func main() {
	var port = flag.Int("p", 8080, "Port to run the server")
	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/", GetHandlerFunc).Methods("GET")
	r.HandleFunc(route, PostHandlerFunc).Methods("POST")

	fmt.Println("Listening on port: ", *port)
    http.ListenAndServe(fmt.Sprintf(":%v", *port), r)
}

func dbConn(ctx context.Context) (*mongo.Client, *mongo.Collection) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("Database connection error")
	}

	recordDatabase := client.Database("Record")
	recordCollection := recordDatabase.Collection("Record")

	return client, recordCollection
}

func GetHandlerFunc(w http.ResponseWriter, r *http.Request){
	// vars := mux.Vars(r)
	// userName := vars["name"]
	// userDate := vars["date"]
	// userTime := vars["time"]

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

	
	var results []*Entry

	for cur.Next(context.Background()){
		var res Entry

		err := cur.Decode(&res)
		if err != nil {
			log.Fatal("Decode error: ",err)
		}
		
		results = append(results, &res)
	}
	fmt.Fprintf(w, "<h1>HELLO WORLD</h1>\n%v", results)
}

func PostHandlerFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["name"]
	userDate := vars["date"]
	userTime := vars["time"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, recordCollection := dbConn(ctx)

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Error: ",err)
		}
	}()

	res, err := recordCollection.InsertOne(ctx, bson.D{
		{"name", userName},
    	{"time", userTime},
		{"date", time.Now()},
		})
	if err != nil {
		log.Fatal("Insert error: ", err)
	}
	fmt.Println(res)
	fmt.Fprintf(w, "<h1>HELLO WORLD</h1>\n%v %v %v", userName, userDate, userTime)
}