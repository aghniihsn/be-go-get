package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	idStr := "687e76c855d32ca936351cfe"

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		fmt.Printf("Error converting ID: %v\n", err)
		return
	}

	fmt.Printf("Valid ObjectID: %s\n", id.Hex())
	fmt.Printf("Timestamp: %v\n", id.Timestamp())
}
