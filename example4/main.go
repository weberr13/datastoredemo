package main

import (
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

type MyNewString struct {
	Number int
	S      string // <-
	K      *datastore.Key `datastore:"__key__"`
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

	for i := 0; i < 100; i++ {
		helloKey := datastore.NameKey("MyNewString", fmt.Sprintf("hello %v", i), nil)
		helloData := MyNewString{}
		helloData.S = fmt.Sprintf("world %v", i)
		helloData.Number = i
		ctx, cancel := context.WithTimeout(pctx, 1*time.Second)
		_, err = cl.Put(ctx, helloKey, &helloData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cancel()
	}
	ctx, cancel := context.WithTimeout(pctx, 1*time.Second)
	q := datastore.NewQuery("MyNewString").Filter("Number =", 25)
	for t := cl.Run(ctx, q); ; {
		var e MyNewString
		_, err := t.Next(&e)
		if err == iterator.Done {
			break
		}
		if e.Number != 25 {
			fmt.Println("query failed")
			os.Exit(1)
		}
		fmt.Println(e.K.Name, e.S)
	}
	cancel()

}
