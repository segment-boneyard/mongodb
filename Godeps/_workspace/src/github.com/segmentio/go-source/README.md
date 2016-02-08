# go-source

The Segment batch data sources library. Use this library to create batch data sources which send object data to the Segment.set API.

# Documentation

[![GoDoc](https://godoc.org/github.com/segmentio/go-source?status.svg)](https://godoc.org/github.com/segmentio/go-source)

# Usage

```go
import source "github.com/segmentio/go-source"

source, err := source.New(&source.Config{
  URL: "http://localhost:4000/rpc",
})
```

# Set API

```go
err := source.Set("leads", "00Q31000013jvx7", map[string]interface{}{
  "first_name": "Bill",
  "last_name": "Lumbergh",
  "email": "bill.lumbergh@initech.com",
})
```
