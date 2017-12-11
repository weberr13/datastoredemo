package util

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
	. "github.com/weberr13/datastoredemo/structs"
)

func RemoveAllKeys(pctx context.Context, cl *datastore.Client, removeKeys []*datastore.Key) (err error) {

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
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
	cancel()

	if err != nil {
		fmt.Println("remove failed: ", err)
		return err
	}
	return nil
}
