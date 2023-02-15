package httptransport

import (
	"bytes"
	"encoding/json"
	"first/chnl"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var templates *template.Template
var client *redis.Client
var store = sessions.NewCookieStore([]byte("t0p-s3cr3t"))

func NewHandler(storage chnl.Store) http.Handler {
	r := mux.NewRouter()

	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "text/html")
			h.ServeHTTP(w, r)
		})
	})

	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	templates = template.Must(template.ParseGlob("templates/*.html"))

	r.HandleFunc("/", indexGetHandler(storage)).Methods("GET")
	r.HandleFunc("/", indexPostHandler(storage)).Methods("POST")
	r.HandleFunc("/login", loginGetHandler(storage)).Methods("GET")
	r.Path("/login").Methods(http.MethodPost).HandlerFunc(loginPostHandler(storage))

	// r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/register", registerGetHandler(storage)).Methods("GET")
	r.HandleFunc("/register", registerPostHandler(storage)).Methods("POST")

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	return r
}

func indexGetHandler(storage chnl.Store) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Print("1\n")

		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]

		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}

		comments, err := client.LRange("comments", 0, 10).Result()

		if err != nil {
			return
		}

		templates.ExecuteTemplate(w, "index.html", comments)
	}
}

func indexPostHandler(storage chnl.Store) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		comment := r.PostForm.Get("comment")
		log.Print(comment)
		client.LPush("comments", comment)
		http.Redirect(w, r, "/", 302)
	}
}

func loginGetHandler(storage chnl.Store) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("2\n")
		templates.ExecuteTemplate(w, "login.html", nil)
	}
}

func loginPostHandler(storage chnl.Store) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		hash, err := client.Get("user:" + username).Bytes()

		if err != nil {
			return
		}

		err = bcrypt.CompareHashAndPassword(hash, []byte(password))

		if err != nil {
			return
		}

		session, _ := store.Get(r, "session")
		session.Values["username"] = username
		session.Save(r, w) // save the session

		//        EDA

		url := "http://localhost:9090/"
		postBody, _ := json.Marshal("login Successful")
		responseBody := bytes.NewBuffer(postBody)
		contentType := "application/json"
		content := responseBody

		s1 := chnl.Collection{url, contentType, content}

		storage.Insert(r.Context(), s1)

		time.Sleep(2 * time.Second)

		http.Redirect(w, r, "/", 302)
	}

}

func registerGetHandler(storage chnl.Store) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "register.html", nil)
	}
}

func registerPostHandler(storage chnl.Store) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		cost := bcrypt.DefaultCost
		hash, err := bcrypt.GenerateFromPassword([]byte(password), cost) // include cost parameter also

		if err != nil {
			return
		}

		client.Set("user:"+username, hash, 0)

		http.Redirect(w, r, "/login", 302)
	}
}

// func loginPostHandler(w http.ResponseWriter, r *http.Request) {

// 	r.ParseForm()
// 	username := r.PostForm.Get("username")
// 	password := r.PostForm.Get("password")

// 	hash, err := client.Get("user:" + username).Bytes()

// 	if err != nil {
// 		return
// 	}

// 	err = bcrypt.CompareHashAndPassword(hash, []byte(password))

// 	if err != nil {
// 		return
// 	}

// 	session, _ := store.Get(r, "session")
// 	session.Values["username"] = username
// 	session.Save(r, w) // save the session

// 	//        EDA

// 	// url := "http://localhost:9090/"
// 	// postBody, _ := json.Marshal("login Successful")
// 	// responseBody := bytes.NewBuffer(postBody)
// 	// contentType := "application/json"
// 	// content := responseBody

// 	// s1 := collection{url, contentType, content}

// 	// s1.MakeChan()

// 	time.Sleep(2 * time.Second)

// 	http.Redirect(w, r, "/", 302)
// }
