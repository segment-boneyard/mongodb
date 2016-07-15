package main

type Column struct {
	Source   string            `json:"source"`
	DataType string            `json:"data_type"`
	Columns  map[string]string `json:"columns"`
}
