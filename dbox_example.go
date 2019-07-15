package main

import (
	"fmt"
	"log"
	"time"
	//"math/rand"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
)

type DataUserModel struct {
	ID        string    `json:"ID" bson:"_id"`
	Name      string    `json:"Name" bson:"Name"`
	Age 	  int		`json:"Age" bson:"Age"`
	Birthday  time.Time	`json:"Birthday" bson:"Birthday"`
	Parents   []string  `json:"Parents" bson:"Parents"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt"`
}

func (*DataUserModel) TableName() string {
	return "datausers"
}

func main() {
	// creating connection
	// need to import github.com/eaciit/dbox/dbc/mongo for connecting to mongodb
	connInfo := dbox.ConnectionInfo{
		Host:     "localhost:27017",
		Database: "belajar_golang",
		UserName: "",
		Password: "",
	}
	conn, err := dbox.NewConnection("mongo", &connInfo)
	if err != nil {
		log.Fatal(err.Error())
	}

	// connect to the database server
	// defer close connection
	err = conn.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	if conn != nil {
		defer conn.Close()
	}

	// =========================================================================================

	// example save operation
	// operasi save akan melakukan insert jika ID belum ada
	// dan akan melakukan update/replace data jika ID sudah ada.
	// data yg akan di insertkan ditaruh dalam Exec(toolkit.M{"data": DATA})
	sample1 := DataUserModel{
		ID:        `rand.Intn(100)`,
		Name:      "noval",
		Age: 		19,
		Birthday:   time.Now(),
		Parents:   []string{"Agus", "Kotak"},
		CreatedAt: time.Now(),
	}
	err = conn.NewQuery().
		From(new(DataUserModel).TableName()).
		Save().
		Exec(toolkit.M{"data": sample1})
	if err != nil {
		log.Fatal(err.Error())
	}

	// =========================================================================================

	// example delete operation
	// delete data dengan query { _id: { $eq: "user2" } } atau _id sama dengan "user2"
	// sisipkan kondisi where pada .Where()
	// gunakan dbox.Eq, dbox.Ne, dbox.Gte, dbox.Lte, dan lainnya untuk operasi query mongo
	err = conn.NewQuery().
		From(new(DataUserModel).TableName()).
		Delete().
		Where(dbox.Eq("_id", "user2")).
		Exec(nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	// =========================================================================================

	// example insert operation
	// jika ID sudah ada sebelumnya, maka akan menghasilkan error
	sample2 := DataUserModel{
		ID:        "user2",
		Name:      "agung",
		//Hobbies:   []string{"eat"},
		CreatedAt: time.Now(),
	}
	err = conn.NewQuery().
		From(new(DataUserModel).TableName()).
		Insert().
		Exec(toolkit.M{"data": sample2})
	if err != nil {
		log.Fatal(err.Error())
	}

	// =========================================================================================

	// example get all operation
	// operasi select sedikit berbeda dengan operasi insert/update
	// tidak menggunakan .Exec() melainkan menggunakan Cursor()
	// lalu dari cursor, akses .Fetch()
	csr1, err := conn.NewQuery().
		From(new(DataUserModel).TableName()).
		Select().
		Cursor(nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	// cursor harus di defer close untuk avoid memory leak!
	defer csr1.Close()

	result1 := make([]DataUserModel, 0)
	err = csr1.Fetch(&result1, 0, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("output1", toolkit.JsonString(result1))

	// =========================================================================================

	// example get by id
	// gunakan dbox.Eq yang merupakan ekuivalen dari $eq di mongo
	csr2, err := conn.NewQuery().
		From(new(DataUserModel).TableName()).
		Select().
		Where(dbox.Eq("_id", "user1")).
		Cursor(nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	// cursor harus di defer close untuk avoid memory leak!
	defer csr2.Close()

	result2 := make([]DataUserModel, 0)
	err = csr2.Fetch(&result2, 0, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("output2", toolkit.JsonString(result2))

	// =========================================================================================

	// example get by 2 kondisi
	// gunakan dbox.And atau dbox.Or untuk operasi perbandingan
	csr3, err := conn.NewQuery().
		From(new(DataUserModel).TableName()).
		Select().
		Where(dbox.Or(
			dbox.Eq("_id", "user1"),
			dbox.Eq("Nama", "agung"),
		)).
		Cursor(nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	// cursor harus di defer close untuk avoid memory leak!
	defer csr3.Close()

	result3 := make([]DataUserModel, 0)
	err = csr3.Fetch(&result3, 0, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("output3", toolkit.JsonString(result3))

	// =========================================================================================

	// aggregate example
	// gunakan .Pipe untuk example aggregate
	pipe := []toolkit.M{

		// gunakan aggergate $match untuk menambahkan WHERE condition
		// pada contoh berikut kita akan filter semua data yg field Hobbies berisi "eat"
		// $elemMatch digunakan disini karena hobbies merupakan array
		toolkit.M{"$match": toolkit.M{
			"Hobbies": toolkit.M{"$elemMatch": toolkit.M{"$eq": "eat"}},
		}},

		// kita group dengan kondisi grouping tanpa kondisi, lalu kita hitung jumlah data yg ada
		toolkit.M{"$group": toolkit.M{
			"_id":   nil,
			"Total": toolkit.M{"$sum": 1},
		}},
	}

	// objek pipe di atas equivalen dengan query aggregate berikut
	// db.datausers.aggregate([
	// 	{ $match: {
	// 		Hobbies: { $elemMatch: { $eq: "eat" } }
	// 	} },
	// 	{ $group: {
	// 		_id: nil,
	// 		Total: { $sum: 1 }
	// 	} }
	// ])

	csr4, err := conn.NewQuery().
		From(new(DataUserModel).TableName()).
		Command("pipe", pipe).
		Cursor(nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	// cursor harus di defer close untuk avoid memory leak!
	defer csr4.Close()

	result4 := make([]toolkit.M, 0)
	err = csr4.Fetch(&result4, 0, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("output4", toolkit.JsonString(result4[0]))
}
