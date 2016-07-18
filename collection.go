package main

type Collection struct {
	Fields         map[string]Field `json:"fields"`
	CollectionName string           `json:"-"`
}

func (c *Collection) GetFieldNames() []string {
	keys := make([]string, len(c.Fields))
	i := 0

	for k := range c.Fields {
		keys[i] = k
		i++
	}

	return keys
}