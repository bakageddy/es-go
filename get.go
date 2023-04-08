package elasticsearch

import (
	"encoding/json"
	"log"
	"os"
	"bytes"
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var client *elasticsearch.Client

func init() {
	var err error
	crt, _ = os.ReadFile("./http_ca.crt")
	cfg := elasticsearch.Config {
		Addresses: []string {
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "dineshkumar",
		CACert: string(crt),
	}

}

func main() {

}
