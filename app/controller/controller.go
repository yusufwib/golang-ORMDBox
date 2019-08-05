package controller

import (
	"encoding/json"
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	. "github.com/eaciit/orm"
	"github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

type DataUserModel struct {
	Context        *orm.DataContext
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
	q, err := conn.NewQuery().
		From("datausers").
		Cursor(nil)
	if err != nil {
		panic("Query Failed");
	}
	users := make([]DataUserModel, 0)
	e = q.Fetch(&users, 0, false)
	if e != nil {
		toolkit.Errorf("Unable to iterate cursor %s", e.Error())
	} else {
		message := fmt.Sprintf(`%s`, toolkit.JsonString(users))
		w.Write([]byte(message))
		return
	}

}
var letterRunes = []rune("QWERTYUIOPnfnffnfnfnASDFGHJKLZXCVBNMADFHIFHWEIOFFBEJBFEOJBFOFBAJFB13124U120286412701123456789")


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
		fmt.Println(payload.Name)

		//connection to database
			ctx, _ := prepareContext()
			defer ctx.Close()
			//ctx.DeleteMany(new(DataUserModel), nil)
		t0 := time.Now()
		count := 1
		b := make([]rune, 10)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		for i := 1; i <= count; i++ {
			fmt.Printf("Insert user no %d ...", i)
			ageNow := time.Now().Year() - payload.Birthday.Year()
			u := new(DataUserModel)
			u.ID = string(b)
			u.Name = payload.Name
			u.Birthday = payload.Birthday
			u.Age = ageNow
			u.Parents = payload.Parents
			u.CreatedAt = time.Now()
			e = ctx.Insert(u)
			fmt.Println(e)
			if e != nil {
				toolkit.Errorf("Error Load %d: %s", i, e.Error())
				return
			} else {
				fmt.Println("OK")
			}
		}
		fmt.Printf("Run process for %v \n", time.Since(t0))

	}

}
func Delete(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		ID string
	}{}

	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, _ := prepareContext()
	defer ctx.Close()
	u := new(DataUserModel)
	e = ctx.GetById(u, payload.ID)
	if e == nil {
		fmt.Printf("Will Delete UserModel:\n %s \n", toolkit.JsonString(u))
		e = ctx.Delete(u)
		if e != nil {
			toolkit.Errorf("Error Load: %s", e.Error())
			return
		}
	} else {
		toolkit.Errorf("Delete error: %s", e.Error())
	}
}

func Update(w http.ResponseWriter, r *http.Request) {

	//inisialisasi
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		payload := struct {
			Name     string   		`json:"Name"`
			Birthday time.Time      `json:"Birthday"`
			Parents  []string 		`json:"Parents"`
			ID 		string
		}{}

		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(payload.Name)

		//connection to database
		ctx, _ := prepareContext()
		defer ctx.Close()
		//ctx.DeleteMany(new(DataUserModel), nil)
		t0 := time.Now()
		count := 2
		b := make([]rune, 10)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		for i := 1; i <= count; i++ {
			fmt.Printf("Insert user no %d ...", i)
			ageNow := time.Now().Year() - payload.Birthday.Year()
			//fmt.Print(ageNow)
			u := new(DataUserModel)
			e = ctx.GetById(u, payload.ID)
			if e == nil {
				fmt.Printf("Will Delete UserModel:\n %s \n", toolkit.JsonString(u))
				e = ctx.Delete(u)
				if e != nil {
					toolkit.Errorf("Error Load: %s", e.Error())
					return
				}
			} else {
				toolkit.Errorf("Delete error: %s", e.Error())
			}
			u.ID = payload.ID
			u.Name = payload.Name
			u.Birthday = payload.Birthday
			u.Age = ageNow
			u.Parents = payload.Parents
			u.CreatedAt = time.Now()
			e = ctx.Insert(u)
			fmt.Println(e)
			if e != nil {
				toolkit.Errorf("Error Load %d: %s", i, e.Error())
				return
			} else {
				fmt.Println("OK")
			}
		}
		fmt.Printf("Run process for %v \n", time.Since(t0))
	}
}
