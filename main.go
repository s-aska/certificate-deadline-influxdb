package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var domains []string
var url string

func main() {
	domains = strings.Split(os.Getenv("DOMAINS"), ",")
	url = os.Getenv("INFLUXDB_WRITE_URL") // http://localhost:8086/write?db=mydb
	port := os.Getenv("PORT")
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.ListenAndServe(":" + port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	for _, domain := range domains {
		t := check(domain)
		duration := t.Unix() - time.Now().Unix()
		post(url, domain, fmt.Sprint(duration))
	}
	fmt.Fprintf(w, "ok")
}

func check(domain string) time.Time {
	config := tls.Config{}

	conn, err := tls.Dial("tcp", domain+":443", &config)
	if err != nil {
		log.Fatal("domain: " + domain + ", error: " + err.Error())
	}

	state := conn.ConnectionState()
	certs := state.PeerCertificates

	defer conn.Close()

	return certs[0].NotAfter
}

func post(url string, domain string, value string) {
	fmt.Println("domain:" + domain + " expires:" + value + " url:" + url)
	client := new(http.Client)
	req, _ := http.NewRequest("POST", url, strings.NewReader("deadline,domain="+domain+" value="+value))
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()
}
