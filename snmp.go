package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	g "github.com/soniah/gosnmp"
)

func main() {

	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = "192.168.10.15"
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0", "1.3.6.1.6.3.1.1.6.1.0"}
	result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	fmt.Println(result.Variables[2].Value)

	// -----------------------------------------------  HTTP PUT part -----------------------------------------------

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // Disable TLS check due to bad certificate

	url := "https://192.168.10.11/_pxc_api/api/variables"
	method := "PUT"
	y := result.Variables[2].Value
	//	x := "{  \"pathPrefix\": \"Arp.Plc.Eclr/\",  \"variables\":  [   {    \"path\": \"bString\",    \"value\": \"%d\",    \"valueType\": \"Constant\"    }    ]}"
	z := fmt.Sprintf("{  \"pathPrefix\": \"Arp.Plc.Eclr/\",  \"variables\":  [   {    \"path\": \"bString\",    \"value\": \"%d\",    \"valueType\": \"Constant\"    }    ]}", y)

	payload := strings.NewReader(z)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	//fmt.Println(x, y, z)
}
