package mymongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func FindMany(coll *mongo.Collection, ctx context.Context, filter interface{}) (interface{}, error) {
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	var all []interface{}
	err = cur.All(context.Background(), &all)
	if err != nil {
		log.Fatal(err)
	}
	cur.Close(context.Background())

	log.Println("collection.Find curl.All: ", all)
	for _, one := range all {
		log.Println(one)
	}

	return all, nil
}
