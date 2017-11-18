/*
package main

import (
	"encoding/json"
	"fmt"
	"time"
	//"github.com/gorilla/context" // To use context
	//"gopkg.in/mgo.v2"            // To interact with mongoDB
	//"gopkg.in/mgo.v2/bson"
	//"log"
	"bytes"
	"io/ioutil"
	"net/http"
)

// for saving room and time for each subject
type Schedule struct {
	//varname type			struct_tag
	//ID      bson.ObjectId `json:"id" bson:"_id"`
	Subject string    `json:"subject" bson:"subject"`
	Room    string    `json:"room" bson:"room"`
	Daytime []Daytime `json:"daytime" bson:"daytime"`
}

//
type Daytime struct {
	//varname type			struct_tag
	Day       int       `json:"day" bson:"day"` // 1~7 : Mon~Sun
	TimeStart time.Time `json:"timestart" bson:"timestart"`
	TimeEnd   time.Time `json:"timeend" bson:"timeend"`
}

// Main
func main() {

	url := "http://localhost:8080/schedules"
	
	schedule := Schedule{
		Subject: "AKE",
		Room:    "7601",
		Daytime: []Daytime{
			{
				Day:       1,
				TimeStart: time.Now(),
				TimeEnd:   time.Now(),
			},
			{
				Day:       5,
				TimeStart: time.Now(),
				TimeEnd:   time.Now(),
			},
		},
	}

	b, err := json.Marshal(schedule)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(schedule)
	//fmt.Println(b)
	fmt.Println(string(b))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
*/
