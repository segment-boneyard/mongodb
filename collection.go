package main

type Collection struct {
	Columns        map[string]Column `json:"columns"`
	CollectionName string            `json:"-"`
}

func (c *Collection) GetColumnNames() []string {
	keys := make([]string, len(c.Columns))
	i := 0

	for k := range c.Columns {
		keys[i] = k
		i++
	}

	return keys
}