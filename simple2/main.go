package main

import (
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

type MyString struct {
	S string
	K *datastore.Key `datastore:"__key__"`
}

func main() {
	project := os.Getenv("DATASTORE_PROJECT_ID")
	if project == "" {
		fmt.Println("must specify datastore environment")
		os.Exit(1)
	}
	pctx, pcancel := context.WithCancel(context.Background())
	defer pcancel()
	cl, err := datastore.NewClient(pctx, project)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cl.Close()
	helloKey := datastore.NameKey("MyString", "hello", nil)
	helloData := MyString{S: "world"}
	ctx, cancel := context.WithTimeout(pctx, 1*time.Second)
	defer cancel()
	helloKey, err = cl.Put(ctx, helloKey, &helloData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	world := &MyString{}
	ctx, cancel = context.WithTimeout(pctx, 1*time.Second)
	defer cancel()
	err = cl.Get(ctx, helloKey, world)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("example:")
	fmt.Println(world.K.Name, world.S)

}
