package main

import (
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	. "github.com/weberr13/datastoredemo/structs"
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
	}
	removeKeys := []*datastore.Key{}

	ctx, cancel = context.WithTimeout(pctx, 1*time.Second)
	q := datastore.NewQuery("MyNewString").Limit(-1).Ancestor(parentKey)
	for t := cl.Run(ctx, q); ; {
		var e MyNewString
		k, err := t.Next(&e)
		if err == iterator.Done {
			break
		}
		removeKeys = append(removeKeys, k)
	}
	cancel()

	ctx, cancel = context.WithTimeout(pctx, 10*time.Second)
	_, err = cl.RunInTransaction(ctx,
		func(t *datastore.Transaction) error {
			err := t.DeleteMulti(removeKeys)
			if err != nil {
				return err
			}
			for i := 0; i < 10; i++ {
				for j := 0; j < len(removeKeys); {
					var e MyNewString
					err = t.Get(removeKeys[j], &e)
					if err != nil {
						removeKeys = append(removeKeys[:j], removeKeys[j+1:]...)
					} else {
						j++
					}
				}
				if len(removeKeys) == 0 {
					break
				} else {
					time.Sleep(10 * time.Millisecond)
				}
			}

			return nil
		})
	if err != nil {
		fmt.Println("remove failed: ", err)
		os.Exit(1)
	}
	cancel()
	ctx, cancel = context.WithTimeout(pctx, 1*time.Second)
	q = datastore.NewQuery("MyNewString").Ancestor(parentKey)
	for t := cl.Run(ctx, q); ; {
		var e MyNewString
		k, err := t.Next(&e)
		if err == iterator.Done {
			break
		}
		fmt.Println("delete failed!!!")
		fmt.Println(k, e)
	}
	cancel()

}
