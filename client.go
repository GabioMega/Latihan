package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/data"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func waitForServer(url string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func extractMsg(r *http.Response) (string, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func get() {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get("http://127.0.0.1:7777/get")
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Println("req timeout")
		} else {
			panic(err)
		}
		return
	}
	defer resp.Body.Close()
	text, err := extractMsg(resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
	return
}
func post() {
	scan := bufio.NewScanner(os.Stdin)
	fmt.Print("title: ")
	scan.Scan()
	title := scan.Text()
	fmt.Print("text: ")
	scan.Scan()
	text2 := scan.Text()

	title = title + ".txt"
	file, err := os.Create(title)
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString("text: " + text2)
	if err != nil {
		panic(err)
	}
	file.Close()
	file, err = os.Open(title)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	multi := multipart.NewWriter(body)
	part, err := multi.CreateFormFile("file", file.Name())
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}
	multi.Close()
	req, err := http.NewRequest("POST", "http://127.0.0.1:7777/post", body)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Println("req timeout")
		} else {
			panic(err)
		}
		return
	}
	req.Header.Set("Content-Type", multi.FormDataContentType())
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	data, err := extractMsg(resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
	return
}

func postjson() {
	var name string
	var age int
	scan := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("input name: ")
		name, _ = scan.ReadString('\n')
		name = strings.TrimSpace(name)

		if len(name) > 5 {
			break
		} else {
			fmt.Println("harus > 5 karakter")
		}
	}
	for {
		fmt.Print("age: ")
		fmt.Scanf("%d\n", &age)
		if age > 0 {
			break
		} else {
			fmt.Println("umur harus > 0 tahun")
		}
	}

	person := data.Person{Name: name, Age: age}
	fmt.Printf("Person struct: %+v\n", person)
	jsonData, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}
	var reqBody bytes.Buffer
	w := multipart.NewWriter(&reqBody)
	personField, err := w.CreateFormField("Person")
	if err != nil {
		panic(err)
	}
	_, err = personField.Write(jsonData)
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}

	//req timeout
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Post("http://127.0.0.1:7777/json", w.FormDataContentType(), &reqBody)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Println("req timeout")
		} else {
			panic(err)
		}
		return
	}
	defer resp.Body.Close()

	data, err := extractMsg(resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
	return
}
func main() {
	scan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("1. get")
		fmt.Println("2. file post")
		fmt.Println("3. json post")
		fmt.Println("4. exit")
		fmt.Print(">> ")
		scan.Scan()
		choice := scan.Text()
		if choice == "1" {
			get()
		} else if choice == "2" {
			post()
		} else if choice == "3" {
			postjson()
		} else if choice == "4" {
			fmt.Println("thanks")
			break
		} else {
			fmt.Println("retype pls")
		}
	}
}
