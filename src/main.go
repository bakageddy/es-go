package main

import (
	// "fmt"
	// "context"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type UploadData struct {
	Id		int64  `json:"id"`
	Data	string `json:"data"`
}

var client *elasticsearch.Client

func init() {
	var err error
	cert, _ := os.ReadFile("./http_ca.crt")

	cfg := elasticsearch.Config {
		Addresses: []string {
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "dineshkumar",
		CACert: cert,
	}
	client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client %s\n", err)
	}
}



func main() {
	readBytes, err := os.ReadFile("./foo.log")

	if err != nil {
		log.Println("Error reading file!\nExiting!\n")
		os.Exit(1)
	}

	readData := string(readBytes)

	document := UploadData {
		Id: 1,
		Data: readData,
	}

	data, err := json.Marshal(document)

	if err != nil {
		log.Println("Error!")
	}

	req := esapi.IndexRequest {
		Index: "test",
		DocumentID: "1",
		Body: bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), client)
	
	if err != nil {
		log.Fatalf("Error %s!\n", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Println("Error!")
		os.Exit(1)
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error!\n");
			os.Exit(1)
		} else {
			log.Printf("%s %s\n", res.Status(), r["result"])
		}
	}

	var r map[string]interface{}

	var buf bytes.Buffer
	query := map[string]interface{} {
		"query" : map[string]interface{} {
			"match" : map[string]interface{} {
				"title": "test",
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalln("Error encoding query!")
		os.Exit(1)
	}

	res, err = client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex("test"),
		client.Search.WithBody(&buf),
		client.Search.WithPretty(),
	)

	if err != nil {
		log.Println("Error!")
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalln("Error!")
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalln("Error parsing response body")
	}

	log.Printf("%v\n", int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)))

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Println(hit)
	}
}
