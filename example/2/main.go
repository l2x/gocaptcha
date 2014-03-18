package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/l2x/gocaptcha"
	"html/template"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/get/", Get)
	http.HandleFunc("/check/", Check)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpl/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func Get(w http.ResponseWriter, r *http.Request) {
	capt := gocaptcha.New()
	f, txt, err := capt.Create()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(txt)
	session, _ := store.Get(r, "session")
	session.Values["session-cap"] = txt
	session.Save(r, w)

	w.Write(f.Bytes())
}

func Check(w http.ResponseWriter, r *http.Request) {
	txt := r.URL.Query().Get("txt")

	session, err := store.Get(r, "session")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err.Error())
		return
	}
	v, ok := session.Values["session-cap"].(string)
	if !ok {
		fmt.Fprint(w, "wrong")
		return
	}

	if v != txt {
		fmt.Fprint(w, v+"!="+txt)
		return
	}

	fmt.Fprint(w, "ok")
}
