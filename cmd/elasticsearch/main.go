package main

import (
	"context"
	"fmt"

	"spacecraft/domain"
	"spacecraft/logic"
)

func main() {
	ctx := context.Background()
	// ctx = logic.ConnectWithElasticSearch(ctx)
	// ctx, err := logic.LoadspacecraftAsync(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	ctx, err := logic.Loadspacecraft(ctx)
	if err != nil {
		panic(err)
	}
	spacecraft := ctx.Value(domain.ModelsKey).([]*domain.Spacecraft)
	if spacecraft == nil {
		panic("no spacecraft in context")
	}
	// if err := internal.WritespacecraftToFile("domain/spacecraft.json", spacecraft); err != nil {
	// 	fmt.Println(fmt.Errorf("WritespacecraftToFile() err: %v", err))
	// }
	// internal.Printspacecraft(os.Stdout, spacecraft)
	fmt.Println("number of spacecraft:", len(spacecraft))
}
