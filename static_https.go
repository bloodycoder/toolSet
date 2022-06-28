package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type fileItem struct{
	Href string
	Title string
}
type Configuration struct{
	RootDir string
}
var rootDir string
func TemplatedHandler(response http.ResponseWriter, request *http.Request) {
	tmplt := template.New("index.html")
	tmplt, _ = tmplt.ParseFiles("index.html")
	fileList := make([]fileItem, 0)
	//read files
	files, err := ioutil.ReadDir(rootDir)
	if err != nil{
		log.Fatal(err)
	}
	for _, file := range files{
		fileList = append(fileList, fileItem{
			Href:  "/"+file.Name(),
			Title: file.Name(),
		})
	}
	tmplt.Execute(response, fileList)
}
func main() {
	conf, _ := os.Open("conf.json")
	defer conf.Close()
	decoder := json.NewDecoder(conf)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatalln(err)
	}
	rootDir = configuration.RootDir
	// start server
	fs := http.FileServer(http.Dir(rootDir))
	http.Handle("/", fs)
	http.HandleFunc("/list", TemplatedHandler)
	log.Print("Listening on :80...")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
