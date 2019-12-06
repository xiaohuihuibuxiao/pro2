package main

import (
	"log"
	"net/http"
)

func main() {
	r := Init()
	log.Println("server starts ...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
