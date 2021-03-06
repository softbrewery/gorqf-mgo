# gorqf-mgo

Go library to add support for mgo in gorqf

[![Build Status](https://travis-ci.org/softbrewery/gorqf-mgo.svg?branch=master)](https://travis-ci.org/softbrewery/gorqf-mgo)
[![codecov](https://codecov.io/gh/softbrewery/gorqf-mgo/branch/master/graph/badge.svg)](https://codecov.io/gh/softbrewery/gorqf-mgo)
[![GitHub license](https://img.shields.io/github/license/softbrewery/gorqf-mgo.svg)](https://github.com/softbrewery/gorqf-mgo/blob/master/LICENSE)

## Install
```shell
$ go get github.com/softbrewery/gorqf
$ go get github.com/softbrewery/gorqf-mgo
```

## Usage

### Basic

```go
// Create rqf parser
parser := rqf.NewParser()

// Json filter comming from rest request
jsonFilter := `
{
    "fields": ["-_id","isbn"],
    "order": ["isbn ASC"],
    "limit": 1,
    "offset": 1,
    "where": {"isbn": "A_ISBN"}
}`

// Parse the filter
filter, err := parser.Parse(jsonFilter)
if err != nil {
    // handle error
}

// Create query and inject where
q := mgoSession.DB("").C("books").Find( MgoWhere(filter) )

// Inject rqf filters (fields/order/limit/offset) in MGO query
MgoAddFilters(q, filter)

// Get data
var data BookList
q.All(&data)
```

### Validation

```go
// Create rqf parser
parser := rqf.NewParser()

// only allow isbn/name to be selected
parser.FieldSchema( joi.String().Allow("isbn", "name") )

// only allow isbn to be ordered (ASC/DESC)
parser.OrderSchema( joi.String().Allow("isbn", "isbn ASC", "isbn DESC") )

// only allow paging between 10 and 100 items
parser.LimitSchema( joi.Int().Min(10).Max(100) )

// Don't allow where condition
parser.WhereSchema( joi.Any().Forbidden() )

// Json filter comming from rest request
jsonFilter := `
{
    "fields": ["isbn", "name"],
    "order": ["isbn ASC"],
    "limit": 25,
    "offset": 1
}`

// Parse the filter
filter, err := parser.Parse(jsonFilter)
if err != nil {
    // handle error
}

// Create query
q := mgoSession.DB("").C("books").Find(nil)

// Inject rqf filters (fields/order/limit/offset) in MGO query
MgoAddFilters(q, filter)

// Get data
var data BookList
q.All(&data)
```

---

## Example

### Json Filter
```json
{
    "fields": ["-_id","isbn","meta.active"],
    "order":["isbn ASC"],
    "limit":1,
    "offset":1
}
```
### Mongo DataSet
```json
[
  {
    "_id": "5ad4e9abce1d82568d3851de",
    "name": "Book1",
    "isbn": "C_ISBN",
    "meta": {
      "active": true,
      "created": "2018-04-16T18:21:31.263Z"
    }
  },
  {
    "_id": "5ad4e9abce1d82568d3851df",
    "name": "Book2",
    "isbn": "A_ISBN",
    "meta": {
      "active": true,
      "created": "2018-04-16T18:21:31.263Z"
    }
  },
  {
    "_id": "5ad4e9abce1d82568d3851e0",
    "name": "Book3",
    "isbn": "B_ISBN",
    "meta": {
      "active": true,
      "created": "2018-04-16T18:21:31.263Z"
    }
  }
]
```

### Result
```json
[
  {
    "isbn": "B_ISBN",
    "meta": {
      "active": true
    }
  }
]
```
