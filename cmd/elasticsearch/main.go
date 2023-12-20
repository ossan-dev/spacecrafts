package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"spacecraft/internal/clients"
	"spacecraft/internal/domain"
	"spacecraft/internal/es"

	"github.com/elastic/go-elasticsearch/v8"
)

const userMenuOptions = `
What do you want to do?

1. Lookup spacecraft by DocumentID.
2. Search spacecraft by status and use UID prefix for relevance.
3. Quit`

// docker-compose up
// go run . -mode=sync -nobulk
func main() {
	ctx := context.Background()
	var err error
	modePtr := flag.String("mode", "async", `specify how to run the program ("sync" or "async"). The latter is the default"`)
	noBulkPtr := flag.Bool("nobulk", false, "specify whether to index data sync or async in elasticsearch")
	flag.Parse()

	esClient, err := setUp(ctx, modePtr, noBulkPtr)
	if err != nil {
		log.Fatalf("failed to set elasticsearch up: %v", err)
	}

	for {
		fmt.Fprintln(os.Stdout, userMenuOptions)
		var userSelection string
		fmt.Fscanln(os.Stdin, &userSelection)
		switch userSelection {
		case "1":
			var documentID string
			fmt.Fprintln(os.Stdout, "type in the DocumentID:")
			fmt.Fscanln(os.Stdin, &documentID)
			spacecraft, err := es.GetByID(esClient, "spacecrafts", documentID)
			if err != nil {
				log.Fatal(err)
			}
			spacecraft.String()
		case "2":
			var uidPrefix string
			fmt.Fprintln(os.Stdout, "type in the UID prefix:")
			fmt.Fscanln(os.Stdin, &uidPrefix)
			var status string
			fmt.Fprintln(os.Stdout, "type in the status:")
			fmt.Fscanln(os.Stdin, &status)
			destroyedSpacecraft, count, err := es.SearchByStatusAndUidPrefix(esClient, "spacecrafts", uidPrefix, status)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("count:", count)
			fmt.Println(len(destroyedSpacecraft))
			for _, v := range destroyedSpacecraft {
				v.String()
			}
		case "3":
			os.Exit(0)
		}
	}
}

func setUp(ctx context.Context, modePtr *string, noBulkPtr *bool) (*elasticsearch.Client, error) {
	var err error
	if modePtr != nil && *modePtr == "sync" {
		client := clients.NewClient("http://localhost:8080", &http.Client{})
		spacecrafts, err := client.Load(ctx)
		if err == nil {
			ctx = domain.WithSpacecrafts(ctx, spacecrafts)
		}

	} else {
		ctx, err = clients.LoadspacecraftAsync(ctx, "http://localhost:8080")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, err
	}
	esClient, err := es.Connect("http://localhost:9200")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, err
	}
	if err = es.DeleteIndex(esClient, "spacecrafts"); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil, err
	}
	if noBulkPtr != nil && *noBulkPtr {
		if err = es.IndexSpacecraftAsDocuments(ctx, esClient); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return nil, err
		}
	} else {
		if err = es.IndexSpacecraftAsDocumentsAsync(ctx, esClient); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return nil, err
		}
	}
	return esClient, nil
}
