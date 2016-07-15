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
}

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
go get -u github.com/segment-sources/mongodb
```

The first step is to initialize your schema. You can do so by running `mongodb` with `--init` flag.
```bash
mongodb --init --write-key=ab-200-1alx91kx --hostname=postgres-test.ksdg31bcms.us-west-2.rds.amazonaws.com --port=5432 --username=segment --password=cndgks8102baajls --database=segment -- sslmode=prefer
```
The init step will store the schema of possible tables that the source can sync in `schema.json`. The query will look for tables across all schemas excluding the ones without a `PRIMARY KEY`.

In the `schema.json` example below, our parser found the collection `films` in the database `public`. The `column` list is used to generate `SELECT` statements, you can filter out some fields that you don't want to sync with Segment by removing them from the list.
```json
{
    "test": {
        "products": {
            "columns": {
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


Segment's Objects API requires a unique identifier in order to properly sync your tables, the `PRIMARY KEY` is used as the identifier. Your tables may also have multiple primary keys, in that case we'll concatenate the values in one string joined with underscores.


### Scan
```bash
mongodb --write-key=ab-200-1alx91kx --hostname=postgres-test.ksdg31bcms.us-west-2.rds.amazonaws.com --port=5432 --username=segment --password=cndgks8102baajls --database=segment --sslmode=prefer
```

### Usage
```
Usage:
  mongodb
    [--debug]
    [--init]
    [--concurrency=<c>]
    --write-key=<segment-write-key>
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
  --concurrency=<c>           Number of concurrent table scans [default: 1]
  --hostname=<hostname>       Database instance hostname
  --port=<port>               Database instance port number
  --password=<password>       Database instance password
  --database=<database>       Database instance name
```
