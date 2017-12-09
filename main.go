package main

import (
	"encoding/json"
	"fmt"
	//"time
	"github.com/gorilla/context" // To use context
	"gopkg.in/mgo.v2"            // To interact with mongoDB
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

// for saving room and time for each subject
/*
type schedule struct {
	//varname type			struct_tag
	ID      bson.ObjectId `json:"id" bson:"_id"`
	Subject string        `json:"subject" bson:"subject"`
	Room    string        `json:"room" bson:"room"`
	Daytime []daytime     `json:"daytime" bson:"daytime"`
}

//
type daytime struct {
	//varname type			struct_tag
	Day       int       `json:"day" bson:"day"` // 1~7 : Mon~Sun
	TimeStart time.Time `json:"timestart" bson:"timestart"`
	TimeEnd   time.Time `json:"timeend" bson:"timeend"`
}
*/

type schedule struct {
	//varname type			struct_tag
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Code       string        `json:"code" bson:"code"`
	Subject    string        `json:"subject" bson:"subject"`
	SKS        string        `json:"sks" bson:"sks"`
	ClassNum   string        `json:"class-num" bson:"class-number"`
	Lecturer   []string      `json:"lecturer" bson:"lecturer"`
	StudentAmt int           `json:"student-amount" bson:"student-amount"`
	Daytime    []daytime     `json:"daytime" bson:"daytime"`
}

//
type daytime struct {
	//varname type			struct_tag
	Day       string `json:"day" bson:"day"` // 1~7 : Mon~Sun
	Room      string `json:"room" bson:"room"`
	TimeStart int    `json:"timestart" bson:"timestart"`
	TimeEnd   int    `json:"timeend" bson:"timeend"`
	Type      string `json:"type" bson:"type"`
}

// to run MongoDB use :
// mongod --dbpath="./db"

// to use curl:
// curl -i -X POST -d "{\"subject\":\"AKE\",\"room\":\"7601\"}"
// http://localhost:8080/schedules

// ADAPTER

/* To write code that can be run before and/or after HTTP requests
 * coming to our API.
 * (Especailly useful for creating connection to MongoDB
 * before request handler run and clean it up after it finish)
 */
type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func withDB(db *mgo.Session) Adapter {
	// return the Adapter
	return func(h http.Handler) http.Handler {
		// the adapter (when called should return a new handler)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// copy the database session
			dbsession := db.Copy()
			defer dbsession.Close() // close up

			// save it in the mux context
			context.Set(r, "database", dbsession)

			// pass execution to the original handler
			h.ServeHTTP(w, r)
		})
	}
}

// CONTROLLER
// handle
func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleRead(w, r)
	case "POST":
		handleInsert(w, r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

// HANDLER
// Read
func handleRead(w http.ResponseWriter, r *http.Request) {
	// get db context (close is handled by adapted handler)
	fmt.Println("Catch Read")

	db := context.Get(r, "database").(*mgo.Session)
	// load the schedules (sliced)
	var schedules []*schedule
	if err := db.DB("scheduleapp").C("schedules").
		Find(nil).Sort("-ID").Limit(100).All(&schedules); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// write it out
	if err := json.NewEncoder(w).Encode(schedules); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Insert
func handleInsert(w http.ResponseWriter, r *http.Request) {
	// get db context (close is handled by adapted handler)
	fmt.Println("Catch Insert")

	db := context.Get(r, "database").(*mgo.Session)
	// decode the request body to 'schedule' struct
	var s schedule

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		fmt.Println("Erorr1")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// give the comment a unique ID
	s.ID = bson.NewObjectId()
	// insert it into the database
	if err := db.DB("scheduleapp").C("schedules").Insert(&s); err != nil {
		fmt.Println("Erorr2")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// redirect to it
	//http.Redirect(w, r, "http://localhost:8080/schedules/"+s.ID.Hex(), http.StatusTemporaryRedirect)
	// loops around because the method is still post
	//http.Redirect(w, r, "http://localhost:8080/schedules", http.StatusTemporaryRedirect)
	//w.Header().Set("Content-Type", "application/json")

	sout := []byte("Your input is accepted")
	w.Write(sout)
}

func updateDB(db *mgo.Session) {
	// How to use:
	// 1. Insert credentials below (nim, usename, password)
	// 2. Uncomment 'updateDB()' in main()
	// 3. 'go run main.go dataFetch.go'

	user := User{}
	user.nim = ""
	user.username = ""
	user.password = ""

	fmt.Println("Please wait. Data is being fetched...")

	scheduleArray := fetch(user)

	for _, s := range scheduleArray {
		if err := db.DB("scheduleapp").C("schedules").Insert(&s); err != nil {
			fmt.Println("Erorr2")
			return
		}
	}
}

// Main
func main() {

	// connect to the database
	fmt.Println("Running")
	db, err := mgo.Dial("localhost") //mgo.Dial returns an mgo.Session
	if err != nil {
		log.Fatal("cannot dial mongo ", err)
	}
	defer db.Close() // clean up when we’re done

	// Adapt our handle function using withDB
	h := Adapt(http.HandlerFunc(handle), withDB(db))

	// ===ACTIVATE THIS IF YOU WANT TO UPDATE DB===
	//updateDB(db)

	// add the handler
	http.Handle("/schedules", context.ClearHandler(h))
	//ClearHandler will clean up any memory used by the context.Set method
	// Hint : path.Base("/id/123") and split
	// start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
