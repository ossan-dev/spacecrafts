package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"spacecraft/internal/clients"
	"spacecraft/internal/domain"
	"spacecraft/logic/elastic"
)

// docker-compose up
// go run . -mode=sync -nobulk
func main() {
	ctx := context.Background()
	var err error
	modePtr := flag.String("mode", "async", `specify how to run the program ("sync" or "async"). The latter is the default"`)
	noBulkPtr := flag.Bool("nobulk", false, "specify whether to index data sync or async in elasticsearch")
	flag.Parse()

	// nit: split line 24 to 52 in its own setUp() function

	if modePtr != nil && *modePtr == "sync" {
		// ctx, err = clients.Loadspacecraft(ctx, "http://localhost:8080")
		client := clients.NewClient("http://localhost:8080", &http.Client{})
		spacecraft, err := client.Load(ctx)
		if err == nil {
			ctx = context.WithValue(ctx, domain.ModelsKey, spacecraft)
		}

	} else {
		ctx, err = clients.LoadspacecraftAsync(ctx, "http://localhost:8080")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	ctx, err = elastic.ConnectWithElasticSearch(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	if err = elastic.DeleteIndex(ctx, "spacecrafts"); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	if noBulkPtr != nil && *noBulkPtr {
		if err = elastic.IndexSpacecraftAsDocuments(ctx); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	} else {
		if err = elastic.IndexSpacecraftAsDocumentsAsync(ctx); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	}

	for {
		// nit: use multi-line strings (with backticks) for long text. Also make the text a constant.
		fmt.Fprintln(os.Stdout, "What do you want to do?\n1. Lookup spacecraft by DocumentID.\n2. Search spacecraft by status and use UID prefix for relevance.\n3. Quit")

		var userSelection string
		fmt.Fscan(os.Stdin, &userSelection) // nit: what about using fmt.Scanln() in main.go ? Same for all others Fscan
		switch userSelection {
		case "1":
			var documentID string
			fmt.Fprintln(os.Stdout, "type in the DocumentID:")
			fmt.Fscan(os.Stdin, &documentID)
			spacecraft, err := elastic.QuerySpacecraftByDocumentID(ctx, "spacecrafts", documentID)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error()) // nit: log.Fatal() ?
				return
			}
			spacecraft.Print(os.Stdout) // major: define String() on spacecraft and simply fmt.Println(spacecraft)
		case "2":
			var uidPrefix string
			fmt.Fprintln(os.Stdout, "type in the UID prefix:")
			fmt.Fscan(os.Stdin, &uidPrefix)
			var status string
			fmt.Fprintln(os.Stdout, "type in the status:")
			fmt.Fscan(os.Stdin, &status)
			destroyedSpacecraft, count, err := elastic.SearchByStatusAndUidPrefix(ctx, "spacecrafts", uidPrefix, status)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error()) // nit: log.Fatal() ?
				return
			}
			fmt.Println("count:", count)
			fmt.Println(len(destroyedSpacecraft))
			for _, v := range destroyedSpacecraft {
				v.Print(os.Stdout)
			}
		case "3":
			os.Exit(0)
		}
		// missing default case
	}
	/******************* Debug ****************************/
	// if err := internal.WritespacecraftToFile("domain/spacecraft.json", spacecraft); err != nil {
	// 	fmt.Println(fmt.Errorf("WritespacecraftToFile() err: %v", err))
	// }
	// internal.Printspacecraft(os.Stdout, spacecraft)
	/******************* End of Debug ******************/
}
