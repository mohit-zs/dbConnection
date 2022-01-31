package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type mysqlConfig struct {
	dbDriver string
	user     string
	password string
	dbName   string
}

type Animal struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

func connectToMySQL(conf mysqlConfig) (*sql.DB, error) {
	db, err := sql.Open(conf.dbDriver, conf.user+":"+conf.password+"@/"+conf.dbName)
	//db, err = sql.Open("mysql", "root:insert_password@tcp(127.0.0.1:3306)/animal")

	if err != nil {
		fmt.Println("an error occurred")
		return nil, err
	}
	return db, nil
}

func main() {
	var err error

	config := mysqlConfig{
		dbDriver: "mysql",
		user:     "root",
		password: "insert_password",
		dbName:   "animal",
	}
	db, err = connectToMySQL(config)

	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/animal/", handler)
	err = http.ListenAndServe(":8000", r)

	if err != nil {
		log.Println(err)
	}

}

func handler(writer http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		values := request.URL.Query()
		if len(values) == 0 {
			GetAnimal(writer, request)
		} else {
			GetAnimalByID(writer, request)
		}
		writer.WriteHeader(http.StatusOK)
	case http.MethodPost:
		PostAnimal(writer, request)
		writer.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		DeleteByID(writer, request)
		writer.WriteHeader(http.StatusOK)
	case http.MethodPut:
		UpdateByID(writer, request)
		writer.WriteHeader(http.StatusOK)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func GetAnimal(w http.ResponseWriter, r *http.Request) {

	var a Animal
	var animals []Animal

	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT  * from animals")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something unexpected happened"))
		return
	}

	for rows.Next() {
		err := rows.Scan(&a.ID, &a.Name, &a.Age)
		if err != nil {
			log.Println(err)
		}
		animals = append(animals, a)
	}

	res, err := json.Marshal(animals)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(res)
}

func GetAnimalByID(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("ID")
	//name := r.URL.Query().Get("Name")

	var a Animal
	var animals []Animal
	//vars := mux.Vars(r)
	//age, _ := strconv.Atoi(vars["age"])
	//age := r.URL.Query().Get("Age")
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT * from animals where id=?", id)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		err := rows.Scan(&a.ID, &a.Name, &a.Age)
		if err != nil {
			log.Println(err)
		}
		animals = append(animals, a)
	}

	res, err := json.Marshal(animals)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(res)
}

//
func PostAnimal(writer http.ResponseWriter, request *http.Request) {
	var a Animal
	body, _ := ioutil.ReadAll(request.Body)
	_ = json.Unmarshal(body, &a)
	_, err := db.Exec("INSERT INTO animals (name,age) VALUES(?,?)", a.Name, a.Age)

	if err != nil {
		log.Println("error:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("something unexpected happened"))
		return
	}
	_, _ = writer.Write([]byte("success"))
}

func DeleteByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("ID")

	w.Header().Set("Content-Type", "application/json")
	_, err := db.Query("DELETE from animals where id=?", id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something unexpected happened"))
		return
	}

	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte("Deletion Successful"))
}

func UpdateByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("ID")

	//w.Header().Set("Content-Type", "application/json")
	var a Animal
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &a)

	_, err := db.Exec("UPDATE animals SET name=?, age=? WHERE id=?", a.Name, a.Age, id)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Update Unsuccessful"))
		return
	}
	w.Write([]byte("Update Successful"))
}
