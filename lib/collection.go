package mongodb

type Field struct {
	FieldName       string `json:"-"`
	DestinationName string `json:"destination_name"`
}

type Collection struct {
	CollectionName  string            `json:"-"`
	DestinationName string            `json:"destination_name,omitempty"`
	Fields          map[string]*Field `json:"fields"`
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
