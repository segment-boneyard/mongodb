package main

import (
	"io"
	"encoding/json"
)

type Description struct {
	// database name -> collection name -> collection
	schemas map[string]map[string]*Collection
}

func NewDescription() *Description {
	return &Description{
		schemas: make(map[string]map[string]*Collection),
	}
}

func NewDescriptionFromReader(r io.Reader) (*Description, error) {
	d := NewDescription()
	if err := json.NewDecoder(r).Decode(&d.schemas); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Description) AddCollection(collectionName string, dbName string) {
	if _, ok := d.schemas[dbName]; !ok {
		d.schemas[dbName] = map[string]*Collection{}
	}
	d.schemas[dbName][collectionName] = &Collection{}
	d.schemas[dbName][collectionName].Fields = make(map[string]Field)
}

func (d *Description) Save(w io.Writer) error {
	b, err := json.MarshalIndent(d.schemas, "", "\t")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (d *Description) Iter() <-chan *Collection {
	out := make(chan *Collection)
	go func() {
		for _, collectionMap := range d.schemas {
			for collectionName, collection := range collectionMap {
				collection.CollectionName = collectionName
				out <- collection
			}
		}
		close(out)
	}()
	return out
}
