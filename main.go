package main

import (
	"log"
	"net/http"
)

func main() {
	//https://www.codercto.com/a/66808.html
	r := Init()
	log.Println("server starts ...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
