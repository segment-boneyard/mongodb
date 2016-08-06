package mongodb

type Collection struct {
	Fields          map[string]string `json:"fields"`
	CollectionName  string            `json:"-"`
	DestinationName string            `json:"destination_name,omitempty"`
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
