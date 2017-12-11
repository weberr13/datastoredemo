package structs

import (
	"cloud.google.com/go/datastore"	
)

type MyNewString struct {
	Number int
	S      string
	K      *datastore.Key `datastore:"__key__"`
}
type MyString struct {
	S string
}
type MyInt struct {
	Number int
	K     *datastore.Key `datastore:"__key__"`
}

type MyParent struct {
	K      *datastore.Key `datastore:"__key__"`	
}