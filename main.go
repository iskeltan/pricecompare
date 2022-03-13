package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"./middleware"
	"./models"
	"./views"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// HeaderMenuItem is blabla
type HeaderMenuItem struct {
	Name string `json:name`
	Url  string `json:url`
}

// IndexData is blabla
type IndexData struct {
	HeaderMenuItems []HeaderMenuItem
	ActiveUser      interface{}
}

func index(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("jello")

	// tmpl := template.Must(template.ParseFiles("templates/index.html"))
	fmt.Print(context.Get(r, "User"))
	headerMenu := HeaderMenuItem{Name: "Register", Url: "/profiles/register"}

	IndexData := IndexData{}
	IndexData.HeaderMenuItems = append(IndexData.HeaderMenuItems, headerMenu)
	IndexData.ActiveUser = context.Get(r, "User")
	IndexJSON, err := json.Marshal(IndexData)

	if err != nil {
		panic(err)
	}
	// tmpl.Execute(w, index_data)
	w.Write(IndexJSON)
}

func detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars["title"])
}

func main() {
	const PORT = "4000"
	// ROUTER NOTES: https://gowebexamples.com/routes-using-gorilla-mux/
	r := mux.NewRouter()

	fmt.Println("[+] server will be ready...")

	r.HandleFunc("/", middleware.JsonView(middleware.Authenticated(middleware.Logging(index))))
	r.HandleFunc("/detail/{title}", middleware.Logging(detail))

	profilesRouter := r.PathPrefix("/profiles").Subrouter()
	profilesRouter.HandleFunc("/login", middleware.JsonView(middleware.Logging(views.LoginUser)))
	profilesRouter.HandleFunc("/register", middleware.JsonView(middleware.Logging(views.RegisterUser)))
	profilesRouter.HandleFunc("/update", middleware.JsonView(middleware.Authenticated(middleware.Logging(views.UpdateUser))))

	// for static files.
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	models.Migrate()
	log.Print("[+] server is ready on :" + PORT)
	http.ListenAndServe(":"+PORT, r)

}
