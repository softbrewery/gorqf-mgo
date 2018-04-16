package rqfmgo

import (
	"fmt"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/softbrewery/gorqf/pkg/rqf"
)

// MgoAddFields will add field filter to the Mongo Query
func MgoAddFields(mq *mgo.Query, q *rqf.Filter) *rqf.Filter {
	mgoSelect := bson.M{}
	for _, field := range q.Fields {
		if strings.HasPrefix(field, "-") {
			mgoSelect[field[1:]] = 0
		} else {
			mgoSelect[field] = 1
		}
	}
	fmt.Println(mgoSelect)
	mq.Select(mgoSelect)
	return q
}

// MgoAddOrder will add order filter to the Mongo Query
func MgoAddOrder(mq *mgo.Query, q *rqf.Filter) *rqf.Filter {
	sort := []string{}
	for _, value := range q.Order {
		strippedValue := strings.Split(value, " ")
		mgoValue := strippedValue[0]
		if len(strippedValue) > 1 && strippedValue[1] == "DESC" {
			mgoValue = "-" + strippedValue[0]
		}
		sort = append(sort, mgoValue)
	}
	mq.Sort(sort...)
	return q
}

// MgoAddLimit will add limit field filter to the Mongo Query
func MgoAddLimit(mq *mgo.Query, q *rqf.Filter) *rqf.Filter {
	mq.Limit(q.Limit)
	return q
}

// MgoAddOffset will add offset filter to the Mongo Query
func MgoAddOffset(mq *mgo.Query, q *rqf.Filter) *rqf.Filter {
	mq.Skip(q.Offset)
	return q
}
