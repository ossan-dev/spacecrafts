package main

import (
	"context"
	"fmt"

	"spacecrafts/domain"
	"spacecrafts/logic"
)

func main() {
	ctx := context.Background()
	// ctx = logic.ConnectWithElasticSearch(ctx)
	// ctx, err := logic.LoadSpacecraftsAsync(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	ctx, err := logic.LoadSpacecrafts(ctx)
	if err != nil {
		panic(err)
	}
	spacecrafts := ctx.Value(domain.ModelsKey).([]*domain.Spacecraft)
	if spacecrafts == nil {
		panic("no spacecrafts in context")
	}
	// if err := internal.WriteSpacecraftsToFile("domain/spacecrafts.json", spacecrafts); err != nil {
	// 	fmt.Println(fmt.Errorf("WriteSpacecraftsToFile() err: %v", err))
	// }
	// internal.PrintSpacecrafts(os.Stdout, spacecrafts)
	fmt.Println("number of spacecrafts:", len(spacecrafts))
}
