package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"spacecraft/logic/elastic"
	"spacecraft/logic/webclient"
)

// docker-compose up
// go run . -mode=sync
func main() {
	ctx := context.Background()
	var err error
	modePtr := flag.String("mode", "mode", `specify how to run the program ("sync" or "async"). The latter is the default"`)
	flag.Parse()
	if modePtr != nil && *modePtr == "sync" {
		ctx, err = webclient.Loadspacecraft(ctx)
	} else {
		ctx, err = webclient.LoadspacecraftAsync(ctx)
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
	/******************* Not use this version ************/
	// if err = logic.IndexSpacecraftAsDocuments(ctx); err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// 	return
	// }
	if err = elastic.IndexSpacecraftAsDocumentsAsync(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	spacecraft, err := elastic.QuerySpacecraftByDocumentID(ctx, "spacecrafts", "1000")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	spacecraft.Print(os.Stdout)
	/******************* Debug ****************************/
	// if err := internal.WritespacecraftToFile("domain/spacecraft.json", spacecraft); err != nil {
	// 	fmt.Println(fmt.Errorf("WritespacecraftToFile() err: %v", err))
	// }
	// internal.Printspacecraft(os.Stdout, spacecraft)
	/******************* End of Debug ******************/
}
