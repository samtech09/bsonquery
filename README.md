# bsonquery
bson query builder for official mongodb driver for golang

## Installation
```
go get github.com/samtech09/bsonquery
```

## Usage
First import package
```
import bq github.com/samtech09/bsonquery
```

Create and build query
```
// filter to query data
filter := bq.Builder().
		And(bq.C().EQ("name", "test2"), bq.C().GT("age", 29)).
		Build()

// Find documents by running filer on collection
cur, err := coll.Find(context.Background(), filter)
...
```

Use `bq.And(...)` for querying one or more fields.

## Supported Operators

### Logical operators
bsonquery Function | equivalent operator
------------------- | -------------------
bq.And | `$and`
bq.Or | `$or`
bq.Not | `$not`
bq.Nor | `$nor`


### Conditional and Other
bsonquery Condition | equivalent operator
------------------- | -------------------
bq.C().EQ | `$eq`
bq.C().GT | `$gt`
bq.C().GTE | `$gte`
bq.C().IN | `$in`
bq.C().LT | `$lt`
bq.C().LTE | `$lte`
bq.C().NE | `$ne`
bq.C().NIN | `$nin`
bq.C().Regex | `$regex`
bq.C().Exist | `$exists`


<br />
 Feedback and suggestions are welcomed.