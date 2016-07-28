# MongoDB Source

Segment source for MongoDB. Syncs your production MongoDB instance with [Segment Objects API](https://github.com/segmentio/objects-go).

### Schema
A `products` collection in the `test` database that looks like this in your production MongoDB instance...

```
{
    "name": "Apple",
    "cost": 1.27,
    "translations": {
        "spanish": "manzana",
        "french" : "pomme"
    }
},
{
    "name": "Pear",
    "cost": 2.01,
    "translations": {
        "spanish": "pera"
    }
}
```

would be queryable in your analytics Redshift or Postgres database like this...

```select * from <source-name>.test_products```

> Redshift

| name  | cost  | translations_spanish  | translations_french |
| ----  |:-----:|:---------------------:|:-------------------:|
| Apple | 1.27  | manzana               | NULL                |
| Pear  | 2.01  | pera                  | pomme               |

Note that the user must explicitly define which fields they want imported from their DB ahead of time. See below for more on how that works.

## Quick Start

### Docker

If you're running docker in production, you can simply run this as a docker image:

```
$ docker run segment/mongodb-source <your-options>
```

### Build and Run
Prerequisites: [Go](https://golang.org/doc/install)

```bash
go get github.com/segment-sources/mongodb
```

The first step is to initialize your schema. You can do so by running `mongodb` with `--init` flag.
```bash
mongodb --hostname=mongo-test.ksd31bacms.us-west-2.rds.amazonaws.com --port=27017 --username=segment --password=cndgks9102baajls --database=segment --sslmode=prefer --init
```
The init step will store the schema of possible collections that the source can sync in `schema.json`. The user should then fill in which fields for each collection should be exported. If no fields for a collection are desired, feel free to remove that particular collection from the JSON entry altogether.

In the `schema.json` example below, our parser found the collection `products` in the database `test`.
```json
{
    "test": {
        "products": {
        }
    }
}
```

Let's say a user wants to export 4 fields: `name`, `cost`, `translations_spanish`, `translations_french` as in the original example of this doc. The JSON should then be:
```json
{
    "test": {
        "products": {
            "fields": {
                "name": {
                    "source": "name"
                },
                "cost": {
                    "source": "cost"
                },
                "translations_spanish": {
                    "source": "translations.spanish"
                },
                "translations_french": {
                    "source": "translations.french"
                }
            }
        }
    }
}
```
`name` and `cost` are first level fields, so their `source` values are simply the field names. The other two fields are nested fields so they need to refer to their nested field names using dot syntax, for example `translations.spanish` and `translations.french`.

Some notes:
* :warning: The warehouse type for a particular field is set the first time data for that field is seen. If subsequent data inserted into the warehouse has a different type than the original type seen, the field value may not cast correctly and loaded into the warehouse properly.
* Currently the only supported MongoDB data types are string, integer, long, double, boolean, date.
* Each object's native `_id_` field is already uploaded by default to Segment and is used as a unique identifier for that object. There is no need to put this field in `schema.json`.


### Scan
To begin exporting fields out of the DB, remove the `--init` flag and add a `--write-key` value:
```bash
mongodb --hostname=mongo-test.ksd31bacms.us-west-2.rds.amazonaws.com --port=27017 --username=segment --password=cndgks9102baajls --database=segment --sslmode=prefer --write-key=ab-200-1alx91kx
```

### Usage
```
Usage:
  mongodb
    [--debug]
    [--init]
    [--concurrency=<c>]
    [--write-key=<segment-write-key>]
    --hostname=<hostname>
    --port=<port>
    --username=<username>
    --password=<password>
    --database=<database>
    [-- <extra-driver-options>...]
  mongodb -h | --help
  mongodb --version

Options:
  -h --help                   Show this screen
  --version                   Show version
  --write-key=<key>           Segment source write key
  --concurrency=<c>           Number of concurrent collection scans [default: 1]
  --hostname=<hostname>       Database instance hostname
  --port=<port>               Database instance port number
  --password=<password>       Database instance password
  --database=<database>       Database instance name
```

