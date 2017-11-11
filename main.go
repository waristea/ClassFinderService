package main

import(
	"fmt"
	"log"
	"net/http"
)

// Main
func main(){
	port := 8080
	
	http.HandleFunc("/", helloWorldHandler)
	
	log.Printf("Server starting on port %v \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello World\n")
}
