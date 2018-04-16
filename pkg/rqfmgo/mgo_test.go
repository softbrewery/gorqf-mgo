package rqfmgo_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/softbrewery/gorqf-mgo/pkg/rqfmgo"
	"github.com/softbrewery/gorqf/pkg/rqf"
)

func connectMgo() *mgo.Session {
	session, err := mgo.Dial("localhost/rqf")

	if err != nil {
		panic(err)
	}

	return session
}

type Meta struct {
	Active  bool       `bson:"active" json:"active,omitempty"`
	Created *time.Time `bson:"created" json:"created,omitempty"`
}
type Book struct {
	ID   bson.ObjectId `bson:"_id" json:"id,omitempty"`
	Name string        `bson:"name" json:"name,omitempty"`
	ISBN string        `bson:"isbn" json:"isbn,omitempty"`
	Meta *Meta         `bson:"meta" json:"meta,omitempty"`
}
type BookList []*Book

// Now ...
func Now() *time.Time {
	time := time.Now()
	return &time
}

var _ = Describe("Mgo", func() {

	BeforeSuite(func() {
		mgoSession := connectMgo().Clone()
		defer mgoSession.Close()

		mgoSession.DB("").DropDatabase()

		err := mgoSession.DB("").C("fields_test").Insert(
			&Book{
				ID:   bson.NewObjectId(),
				Name: "Book1",
				ISBN: "C_ISBN",
				Meta: &Meta{
					Active:  true,
					Created: Now(),
				},
			},
			&Book{
				ID:   bson.NewObjectId(),
				Name: "Book2",
				ISBN: "A_ISBN",
				Meta: &Meta{
					Active:  true,
					Created: Now(),
				},
			},
			&Book{
				ID:   bson.NewObjectId(),
				Name: "Book3",
				ISBN: "B_ISBN",
				Meta: &Meta{
					Active:  true,
					Created: Now(),
				},
			})

		if err != nil {
			panic(err)
		}
	})

	Describe("MgoAddFields", func() {

		It("should return all fields if not set", func() {
			mgoSession := connectMgo().Clone()
			defer mgoSession.Close()

			var data BookList

			q := mgoSession.DB("").C("fields_test").Find(nil)
			q.All(&data)

			Expect(len(data)).To(Equal(3))
			for _, item := range data {
				Expect(item.ID).ToNot(BeEmpty())
				Expect(item.Name).ToNot(BeEmpty())
				Expect(item.ISBN).ToNot(BeEmpty())
			}

			fmt.Fprintln(GinkgoWriter, "#####################################################")
			fmt.Fprintln(GinkgoWriter, "No fields filter")
			json, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(GinkgoWriter, string(json))
		})

		It("should return selected fields if set", func() {
			parser := rqf.NewParser()

			jsonFilter := `{"fields":["name"]}`
			filter, err := parser.Parse(jsonFilter)

			Expect(err).To(BeNil())

			mgoSession := connectMgo().Clone()
			defer mgoSession.Close()

			var data BookList

			q := mgoSession.DB("").C("fields_test").Find(nil)
			MgoAddFields(q, filter)
			q.All(&data)

			Expect(len(data)).To(Equal(3))
			for _, item := range data {
				Expect(item.ID).ToNot(BeEmpty())
				Expect(item.Name).ToNot(BeEmpty())
				Expect(item.ISBN).To(BeEmpty())
			}

			fmt.Fprintln(GinkgoWriter, "#####################################################")
			fmt.Fprintln(GinkgoWriter, "Filter: ", jsonFilter)
			json, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(GinkgoWriter, string(json))
		})

		It("should return selected fields if set (-)", func() {
			parser := rqf.NewParser()

			jsonFilter := `{"fields":["-_id", "name"]}`
			filter, err := parser.Parse(jsonFilter)

			Expect(err).To(BeNil())

			mgoSession := connectMgo().Clone()
			defer mgoSession.Close()

			var data BookList

			q := mgoSession.DB("").C("fields_test").Find(nil)
			MgoAddFields(q, filter)
			q.All(&data)

			Expect(len(data)).To(Equal(3))
			for _, item := range data {
				Expect(item.ID).To(BeEmpty())
				Expect(item.Name).ToNot(BeEmpty())
				Expect(item.ISBN).To(BeEmpty())
			}

			fmt.Fprintln(GinkgoWriter, "#####################################################")
			fmt.Fprintln(GinkgoWriter, "Filter: ", jsonFilter)
			json, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(GinkgoWriter, string(json))
		})

		It("Should return fields (sub document filter)", func() {
			parser := rqf.NewParser()

			jsonFilter := `{"fields":["meta.active"]}`
			filter, err := parser.Parse(jsonFilter)

			Expect(err).To(BeNil())

			mgoSession := connectMgo().Clone()
			defer mgoSession.Close()

			var data BookList

			q := mgoSession.DB("").C("fields_test").Find(nil)
			MgoAddFields(q, filter)
			q.All(&data)

			Expect(len(data)).To(Equal(3))
			for _, item := range data {
				Expect(item.ID).ToNot(BeEmpty())
				Expect(item.Name).To(BeEmpty())
				Expect(item.ISBN).To(BeEmpty())
				Expect(item.Meta).ToNot(BeNil())
			}

			fmt.Fprintln(GinkgoWriter, "#####################################################")
			fmt.Fprintln(GinkgoWriter, "Filter: ", jsonFilter)
			json, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(GinkgoWriter, string(json))
		})
	})

	Describe("MgoAddOrder", func() {

		It("Should order results (order isbn)", func() {
			parser := rqf.NewParser()

			jsonFilter := `{"order":["isbn"]}`
			filter, err := parser.Parse(jsonFilter)

			Expect(err).To(BeNil())

			mgoSession := connectMgo().Clone()
			defer mgoSession.Close()

			var data BookList

			q := mgoSession.DB("").C("fields_test").Find(nil)
			MgoAddOrder(q, filter)
			q.All(&data)

			Expect(data[0].ISBN).To(Equal("A_ISBN"))
			Expect(data[1].ISBN).To(Equal("B_ISBN"))
			Expect(data[2].ISBN).To(Equal("C_ISBN"))

			fmt.Fprintln(GinkgoWriter, "#####################################################")
			fmt.Fprintln(GinkgoWriter, "Filter: ", jsonFilter)
			json, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(GinkgoWriter, string(json))
		})

		It("Should order results (order isbn ASC)", func() {
			parser := rqf.NewParser()

			jsonFilter := `{"order":["isbn ASC"]}`
			filter, err := parser.Parse(jsonFilter)

			Expect(err).To(BeNil())

			mgoSession := connectMgo().Clone()
			defer mgoSession.Close()

			var data BookList

			q := mgoSession.DB("").C("fields_test").Find(nil)
			MgoAddOrder(q, filter)
			q.All(&data)

			Expect(data[0].ISBN).To(Equal("A_ISBN"))
			Expect(data[1].ISBN).To(Equal("B_ISBN"))
			Expect(data[2].ISBN).To(Equal("C_ISBN"))

			fmt.Fprintln(GinkgoWriter, "#####################################################")
			fmt.Fprintln(GinkgoWriter, "Filter: ", jsonFilter)
			json, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(GinkgoWriter, string(json))
		})

		It("Should order results (order isbn DESC)", func() {
			parser := rqf.NewParser()

			jsonFilter := `{"order":["isbn DESC"]}`
			filter, err := parser.Parse(jsonFilter)

			Expect(err).To(BeNil())

			mgoSession := connectMgo().Clone()
			defer mgoSession.Close()

			var data BookList

			q := mgoSession.DB("").C("fields_test").Find(nil)
			MgoAddOrder(q, filter)
			q.All(&data)

			Expect(data[0].ISBN).To(Equal("C_ISBN"))
			Expect(data[1].ISBN).To(Equal("B_ISBN"))
			Expect(data[2].ISBN).To(Equal("A_ISBN"))

			fmt.Fprintln(GinkgoWriter, "#####################################################")
			fmt.Fprintln(GinkgoWriter, "Filter: ", jsonFilter)
			json, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(GinkgoWriter, string(json))
		})
	})

	Describe("MgoLimit", func() {

		It("Should limit results (limit 2)", func() {
			parser := rqf.NewParser()

			jsonFilter := `{"limit": 2}`
			filter, err := parser.Parse(jsonFilter)

			Expect(err).To(BeNil())

			mgoSession := connectMgo().Clone()
			defer mgoSession.Close()

			var data BookList

			q := mgoSession.DB("").C("fields_test").Find(nil)
			MgoAddLimit(q, filter)
			q.All(&data)

			Expect(len(data)).To(Equal(2))
			Expect(data[0].Name).To(Equal("Book1"))
			Expect(data[1].Name).To(Equal("Book2"))

			fmt.Fprintln(GinkgoWriter, "#####################################################")
			fmt.Fprintln(GinkgoWriter, "Filter: ", jsonFilter)
			json, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(GinkgoWriter, string(json))
		})
	})

	Describe("MgoOffset", func() {

		It("Should skip results (offset 2)", func() {
			parser := rqf.NewParser()

			jsonFilter := `{"offset": 2}`
			filter, err := parser.Parse(jsonFilter)

			Expect(err).To(BeNil())

			mgoSession := connectMgo().Clone()
			defer mgoSession.Close()

			var data BookList

			q := mgoSession.DB("").C("fields_test").Find(nil)
			MgoAddOffset(q, filter)
			q.All(&data)

			Expect(len(data)).To(Equal(1))
			Expect(data[0].Name).To(Equal("Book3"))

			fmt.Fprintln(GinkgoWriter, "#####################################################")
			fmt.Fprintln(GinkgoWriter, "Filter: ", jsonFilter)
			json, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(GinkgoWriter, string(json))
		})
	})

	Describe("MgoAddFilters", func() {

		It("Should add fields/order/limit/offset filters", func() {
			parser := rqf.NewParser()

			jsonFilter := `{"fields":["-_id","isbn"],"order":["isbn ASC"],"limit":1,"offset":1}`
			filter, err := parser.Parse(jsonFilter)

			Expect(err).To(BeNil())

			mgoSession := connectMgo().Clone()
			defer mgoSession.Close()

			var data BookList

			q := mgoSession.DB("").C("fields_test").Find(nil)
			MgoAddFilters(q, filter)
			q.All(&data)

			Expect(len(data)).To(Equal(1))
			Expect(data[0].ID).To(BeEmpty())
			Expect(data[0].ISBN).To(Equal("B_ISBN"))

			fmt.Fprintln(GinkgoWriter, "#####################################################")
			fmt.Fprintln(GinkgoWriter, "Filter: ", jsonFilter)
			json, _ := json.MarshalIndent(data, "", "  ")
			fmt.Fprintln(GinkgoWriter, string(json))
		})
	})
})
