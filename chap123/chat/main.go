package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/mihirat/go-webapp/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

const (
	GomniSecurityKey = "securty key"
	GoogleClientID   = "666323378036-9jshbanif7srjpoa0n5isgufoudb05ib.apps.googleusercontent.com"
	GoogleSecretKey  = "dADjithikIeH9vjjcRQngR4h"
	AuthCookieName   = "auth"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatarAvatar,
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie(AuthCookieName); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "app address")
	flag.Parse()

	// set up gomniauth
	gomniauth.SetSecurityKey(GomniSecurityKey)
	gomniauth.WithProviders(
		google.New(GoogleClientID, GoogleSecretKey, "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))

	http.Handle("/room", r)
	go r.run()

	log.Println("start web server port: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Listenandserve", err)
	}
}
