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
	fmt.Println("[Startup] domains:"+strings.Join(domains, ",")+" url:"+url)
	checkAll()
	cron()

	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "domains:"+strings.Join(domains, ",")+" url:"+url)
}

func cron() {
	go func() {
		t := time.NewTicker(600 * time.Second)
		for {
			select {
			case <-t.C:
				checkAll()
			}
		}
		t.Stop()
	}()
}

func checkAll() {
	for _, domain := range domains {
		checkDeadline(domain)
		checkElapsed(domain)
	}
}

func checkElapsed(domain string) {
	s := time.Now()
	res, err := http.Get("https://"+domain)
	if err != nil {
		log.Fatal("error:" + err.Error())
		return
	}
	elapsed := time.Since(s)
	defer res.Body.Close()
	reportElapsed(domain, fmt.Sprint(elapsed.Nanoseconds()))
}

func reportElapsed(domain string, value string) {
	fmt.Println("domain:" + domain + " elapsed:" + value + " url:" + url)
	client := new(http.Client)
	req, _ := http.NewRequest("POST", url, strings.NewReader("elapsed,domain="+domain+" value="+value))
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("error:" + err.Error())
		return
	}
	defer res.Body.Close()
}

func checkDeadline(domain string) {
	config := tls.Config{}

	conn, err := tls.Dial("tcp", domain+":443", &config)
	if err != nil {
		log.Fatal("domain:" + domain + " error:" + err.Error())
		return
	}

	state := conn.ConnectionState()
	certs := state.PeerCertificates

	defer conn.Close()

	duration := certs[0].NotAfter.Unix() - time.Now().Unix()
	reportDeadline(domain, fmt.Sprint(duration))
}

func reportDeadline(domain string, value string) {
	fmt.Println("domain:" + domain + " expires:" + value + " url:" + url)
	client := new(http.Client)
	req, _ := http.NewRequest("POST", url, strings.NewReader("deadline,domain="+domain+" value="+value))
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("error:" + err.Error())
		return
	}
	defer res.Body.Close()
}
