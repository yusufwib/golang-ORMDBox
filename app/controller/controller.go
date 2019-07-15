package controller

import (
	"encoding/json"
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	. "github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

type DataUserModel struct {
	ID        string    `json:"ID" bson:"_id"`
	Name      string    `json:"Name" bson:"Name"`
	Age 	  int		`json:"Age" bson:"Age"`
	Birthday  time.Time `json:"Birthday" bson:"Birthday"`
	Parents   []string  `json:"Parents" bson:"Parents"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt"`
}

func (u *DataUserModel) PreSave() error {
	panic("implement me")
}

func (u *DataUserModel) PostSave() error {
	panic("implement me")
}

func (u *DataUserModel) SetM(IModel) IModel {
	panic("implement me")
}

func (u *DataUserModel) PrepareID() interface{} {
	panic("implement me")
}

var e error
func (u *DataUserModel) Init() *DataUserModel {
	//u.M = u
	return u
}
func prepareContext() (*DataContext, error) {
	conn, _ := dbox.NewConnection("mongo", &dbox.ConnectionInfo{"localhost:27017", "belajar_golang", "", "", nil})
	if eConnect := conn.Connect(); eConnect != nil {
		return nil, eConnect
	}
	ctx := New(conn)
	return ctx, nil
}
func (*DataUserModel) TableName() string {
	return "datausers"
}
func (u *DataUserModel) RecordID() interface{} {
	return u.ID
}

var ctx *DataContext


func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("app/views/view.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ci := dbox.ConnectionInfo {
		"127.0.0.1:27017",
		"belajar_golang",
		"",
		"",
		nil,
	}
	conn, err := dbox.NewConnection("mongo", &ci)
	if err != nil {
		panic("Connect Failed"); // Change with your error handling
	}
	err = conn.Connect()
	if err != nil {
		panic("Connect Failed"); // Change with your error handling
	}
	//fmt.Println("Test Load All")
	q, err := conn.NewQuery().
		//Select("Name").
		From("datausers").
		//Where(dbox.Eq("_id", "0123456789")).
		Cursor(nil)
	if err != nil {
		panic("Query Failed");
	}
	// This map is only for quick example
	// It's better to use some struct to create stronger type check
	users := make([]DataUserModel, 0)
	e = q.Fetch(&users, 0, false)
	if e != nil {
		toolkit.Errorf("Unable to iterate cursor %s", e.Error())
	} else {
		//fmt.Print(toolkit.JsonString(users))
		//var jsonObj = json.Marshal()
		//
		//var message = fmt.Sprintf(`%s\n`, toolkit.JsonString(users))
		//var jsonData = []byte(message)
		//var jsonObj, _ = json.Marshal(jsonData)
		//var data DataUserModel
		//
		//var err = json.Unmarshal(jsonObj, &data)
		//if err != nil {
		//	fmt.Println(err.Error())
		//	return
		//}
		//
		//fmt.Println("user :", data.Name)
		//fmt.Println("age  :", data.Age)
		//

		message := fmt.Sprintf(`%s`, toolkit.JsonString(users))
		w.Write([]byte(message))
		return

		//fmt.Printf("\n%s\n", toolkit.JsonString(users))
	}

}
var letterRunes = []rune("qwertyuiopasdfghjklzxcvbnm1234567890")

func CreateUsers(w http.ResponseWriter, r *http.Request) {
	//inisialisasi
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		payload := struct {
			Name     string   		`json:"Name"`
			Birthday time.Time      `json:"Birthday"`
			Parents  []string 		`json:"Parents"`
		}{}
		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//connection to database
		ctx, _ := prepareContext()
		defer ctx.Close()
		//ctx.DeleteMany(new(DataUserModel), nil)
		t0 := time.Now()
		count := 2
		b := make([]rune, 20)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		for i := 1; i <= count; i++ {
			fmt.Printf("Insert user no %d ...", i)
			ageNow := time.Now().Year() - payload.Birthday.Year()
			fmt.Print(ageNow)
			u := new(DataUserModel)
			u.ID = string(b)
			u.Name = payload.Name
			u.Birthday = payload.Birthday
			u.Age = ageNow
			u.Parents = payload.Parents
			e = ctx.Insert(u)
			if e != nil {
				//t.Errorf("Error Load %d: %s", i, e.Error())
				return
			} else {
				fmt.Println("OK")
			}
		}
		fmt.Printf("Run process for %v \n", time.Since(t0))

	}

}