package main

type Field struct {
	Source   string            `json:"source"`
	DataType string            `json:"data_type"`
	Fields   map[string]string `json:"fields"`
}
