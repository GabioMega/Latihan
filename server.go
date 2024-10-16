package main

import (
	"encoding/json"
	"fmt"
	"io"
	"main/data"
	"net/http"
	"os"
	"path/filepath"
)

func middlewareValidation(next http.Handler, accMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range accMethods {
			if method == r.Method {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "method is not allowed. pls recheck", http.StatusBadRequest)
		return
	})
}

func sendJsonHandler(w http.ResponseWriter, r *http.Request) {
	jsonData := r.FormValue("Person")
	var person data.Person
	err := json.Unmarshal([]byte(jsonData), &person)
	if err != nil {
		panic(err)
	}
	fmt.Println("received json data: ", person)

	//response
	fmt.Fprintln(w, "Server got the json")
}

func postResp(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	dst, err := os.Create(filepath.Join("uploads", handler.Filename))
	if err != nil {
		panic(err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Server got the message")
}
func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Good Morning!")
}
func main() {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()

	mux.Handle("/get", middlewareValidation(http.HandlerFunc(handleGet), "GET"))
	mux.Handle("/post", middlewareValidation(http.HandlerFunc(postResp), "POST"))
	mux.Handle("/json", middlewareValidation(http.HandlerFunc(sendJsonHandler), "POST"))

	// err := http.ListenAndServe("127.0.0.1:7777", mux)
	// if err != nil {
	// 	panic(err)
	// }
	server := http.Server{
		Addr:    "127.0.0.1:7777",
		Handler: mux,
	}
	server.ListenAndServe()
}
