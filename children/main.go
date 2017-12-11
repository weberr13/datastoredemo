package main

import (
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	. "github.com/weberr13/datastoredemo/structs"
	"github.com/weberr13/datastoredemo/util"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

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

	parentKey := datastore.NameKey("MyParent", "example8", nil)
	parentData := MyParent{}
	ctx, cancel := context.WithTimeout(pctx, 1*time.Second)
	_, err = cl.Put(ctx, parentKey, &parentData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cancel()
	removeKeys := []*datastore.Key{}

	for i := 0; i < 100; i++ {
		helloKey := datastore.NameKey("MyNewString", fmt.Sprintf("hello %v", i), parentKey)
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
		secretKey := datastore.NameKey("MyInt", "secret", helloKey)
		secretData := MyInt{Number: i + 100}
		ctx, cancel = context.WithTimeout(pctx, 1*time.Second)
		_, err = cl.Put(ctx, secretKey, &secretData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cancel()
		removeKeys = append(removeKeys, helloKey)
	}
	defer func() {
		err = util.RemoveAllKeys(pctx, cl, removeKeys)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	ctx, cancel = context.WithTimeout(pctx, 1*time.Second)
	q := datastore.NewQuery("MyInt").Ancestor(parentKey).Filter("Number =", 125).KeysOnly()
	for t := cl.Run(ctx, q); ; {
		var e MyNewString
		k, err := t.Next(nil)
		if err == iterator.Done {
			break
		}
		worldKey := k.Parent
		err = cl.Get(ctx, worldKey, &e)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Found: ", e)
	}
	cancel()

}
