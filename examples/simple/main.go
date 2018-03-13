package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/tampajohn/snowflake-ingest-go/ingestion"
)

func main() {
	pkeyBytes, err := ioutil.ReadFile("test.pem")
	if err != nil {
		log.Fatal("Must have a private key in pem format in the root directory of this repo in order to run this example.")
	}
	if err != nil {
		log.Fatal("Error reading pem file.")
	}
	m, err := ingestion.NewManager("YOUR_ACCUNT", "YOUR_USER", "YOUR_DB.YOUR_SCHEMA.YOUR_PIPE", &pkeyBytes)
	if err != nil {
		log.Fatal("Invalid arguments passed to NewManager:", err)
	}
	rid := "testid1"
	req := m.NewRequest(&rid)
	resp, err := req.
		AddFiles("somefile.csv.gz").
		DoIngest(context.Background())

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.StatusCode)
		s, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(s))
	}
	rid = "testid2"
	req2 := m.NewRequest(&rid)
	resp2, err := req2.DoReport()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp2.StatusCode)
		s, _ := ioutil.ReadAll(resp2.Body)
		fmt.Println(string(s))
	}

}
