# gorqf-mgo

Go library to add support for mgo in gorqf

[![Build Status](https://travis-ci.org/softbrewery/gorqf-mgo.svg?branch=master)](https://travis-ci.org/softbrewery/gorqf-mgo)
[![codecov](https://codecov.io/gh/softbrewery/gorqf-mgo/branch/master/graph/badge.svg)](https://codecov.io/gh/softbrewery/gorqf-mgo)

## Install
```shell
$ go get github.com/softbrewery/gorqf
$ go get github.com/softbrewery/gorqf-mgo
```

## Usage

```go
parser := rqf.NewParser()

// json filter comming from rest request
jsonFilter := `
{
    "fields": ["-_id","isbn"],
    "order":["isbn ASC"],
    "limit":1,
    "offset":1
}`

// parse the filter
filter, err := parser.Parse(jsonFilter)
if err != nil {
    // return http 400 - bad request
}

// connect to mongo
mgoSession := getMyMongoConnection().Clone()
defer mgoSession.Close()

// create query
q := mgoSession.DB("").C("fields_test").Find(nil)

// inject rqf filters (fields/order/limit/offset) in MGO query
MgoAddFilters(q, filter)

// get data
var data BookList
q.All(&data)
```