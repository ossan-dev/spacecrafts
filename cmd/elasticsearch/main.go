package main

import (
	"context"
	"flag"

	"spacecraft/domain"
	"spacecraft/logic"
)

// docker run -p 9200:9200 -it -m 1GB docker.elastic.co/elasticsearch/elasticsearch:8.11.1
// go run . -mode=sync
func main() {
	ctx := context.Background()
	var err error
	modePtr := flag.String("mode", "mode", `specify how to run the program ("sync" or "async"). The latter is the default"`)
	flag.Parse()
	if modePtr != nil && *modePtr == "sync" {
		ctx, err = logic.Loadspacecraft(ctx)
	} else {
		ctx, err = logic.LoadspacecraftAsync(ctx)
	}
	if err != nil {
		panic(err)
	}
	spacecraft := ctx.Value(domain.ModelsKey).([]*domain.Spacecraft)
	if spacecraft == nil {
		panic("no spacecraft in context")
	}
	ctx = logic.ConnectWithElasticSearch(ctx)
	if err != nil {
		panic(err)
	}
	_ = ctx
	/******************* Debug ****************************/
	// if err := internal.WritespacecraftToFile("domain/spacecraft.json", spacecraft); err != nil {
	// 	fmt.Println(fmt.Errorf("WritespacecraftToFile() err: %v", err))
	// }
	// internal.Printspacecraft(os.Stdout, spacecraft)
	/******************* End of Debug ******************/
}
