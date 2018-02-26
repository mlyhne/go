package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

type Subprogram struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Table struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Program struct {
	Name       string       `json:"name"`
	Version    string       `json:"version"`
	DateCreate string       `json:"date"`
	Calls      []Subprogram `json:"calls"`
	Tables     []Table      `json:"tables"`
}

var tpl *template.Template
var db *sql.DB
var err error

func init() {
	fmt.Printf("Start")
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func GenCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	program := ps.ByName("id")

	// Open our jsonFile
	jsonFile, err := os.Open(program + ".json")
	// if we os.Open returns an error then handle it
	if err != nil {
		io.WriteString(w, "Did not find: "+program+".json (NB case sencitive)")
		return
	}

	fmt.Println("Successfully Opened " + program + ".json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var prog Program

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &prog)

	tpl = template.Must(template.ParseGlob("templates/*"))

	tpl.ExecuteTemplate(w, "gen.gohtml", prog) //only send over images
}

func main() {

	router := httprouter.New()
	router.ServeFiles("/public/*filepath", http.Dir("./public"))
	router.GET("/generate/:id", GenCode)
	router.GET("/generate", GenCode)

	log.Fatal(http.ListenAndServe(":8010", router))
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
