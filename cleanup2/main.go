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
	S      string
	K      *datastore.Key `datastore:"__key__"`
}
type MyString struct {
	S string
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
	removeKeys := []*datastore.Key{}

	ctx, cancel := context.WithTimeout(pctx, 1*time.Second)
	q := datastore.NewQuery("MyNewString").Limit(-1)
	for t := cl.Run(ctx, q); ; {
		var e MyNewString
		k, err := t.Next(&e)
		if err == iterator.Done {
			break
		}
		removeKeys = append(removeKeys, k)
	}
	cancel()
	ctx, cancel = context.WithTimeout(pctx, 1*time.Second)
	q = datastore.NewQuery("MyString").Limit(-1)
	for t := cl.Run(ctx, q); ; {
		var e MyString
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
			for i := 0; i < 10; i++ {
				err := cl.DeleteMulti(ctx, removeKeys)
				if err != nil {
					return err
				}
				q = datastore.NewQuery("MyNewString").KeysOnly()
				_, err = cl.Run(ctx, q).Next(nil)
				if err != iterator.Done {
					continue
				}
				q = datastore.NewQuery("MyString").KeysOnly()
				_, err = cl.Run(ctx, q).Next(nil)
				if err != iterator.Done {
					continue
				}
				break
			}
			q = datastore.NewQuery("MyNewString").KeysOnly()
			_, err = cl.Run(ctx, q).Next(nil)
			if err != iterator.Done {
				return fmt.Errorf("Failed to delete MyNewString: %v", err)
			}
			q = datastore.NewQuery("MyString").KeysOnly()
			_, err = cl.Run(ctx, q).Next(nil)
			if err != iterator.Done {
				return fmt.Errorf("Failed to delete MyString: %v", err)
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cancel()
	ctx, cancel = context.WithTimeout(pctx, 1*time.Second)
	q = datastore.NewQuery("MyNewString")
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
	ctx, cancel = context.WithTimeout(pctx, 1*time.Second)
	q = datastore.NewQuery("MyString").Limit(-1)
	for t := cl.Run(ctx, q); ; {
		var e MyString
		k, err := t.Next(&e)
		if err == iterator.Done {
			break
		}
		fmt.Println("delete failed!!!")
		fmt.Println(k, e)
	}
	cancel()
}
