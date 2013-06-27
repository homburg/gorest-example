package main

import (
	"bytes"
	"code.google.com/p/gorest"
	"github.com/homburg/amber"
	"fmt"
	"net/http"
	"html/template"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var address string

var serviceRoot = "/service/"

var startTime = time.Now()

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	address = "localhost:8787"

	gorest.RegisterService(new(UserService))
	http.Handle(serviceRoot, gorest.Handle())
	http.HandleFunc("/", index)
	http.Handle("/js/", http.FileServer(http.Dir("static")))
	http.Handle("/css/", http.FileServer(http.Dir("static")))

	// Mark file for livereload
	cmd := exec.Command("touch", "server.run")
	cmd.Run()

	http.ListenAndServe(address, nil)
}

var html string

var amberTempl = `doctype 5
html
head
	title Rest exercise

	link[rel="stylesheet"][href="/css/bootstrap.min.css"]
	script[type="text/javascript"][src="/js/components/angular/angular.js"]
	script[type="text/javascript"][src="/js/components/angular-resource/angular-resource.js"]
	script[type="text/javascript"][src="/js/app.js"]
	script[type="text/javascript"][src="http://localhost:35729/livereload.js"]

body[ng-app="myApp"]
	div.container[ng-controller="UsersCtrl"]
		h1 {{ title }}

		p #{Wd}

		input[type="text"][ng-model="title"]

		form[ng-submit="AddUser()"]
			input[type="text"][ng-model="name"][placeholder="navn"]
			input[type="text"][ng-model="email"][placeholder="email"]
			input[type="submit"]
		
		ng-user-details[user="activeUser"]
		textarea {{ activeUser }}

		ul[ng-repeat="user in users"]
			li "{{ user.Name }}" <{{ user.Email }}>
				button[ng-click="ViewUser(user)"] Vis bruger

`

type templateData struct {
	Wd string
}

func index(w http.ResponseWriter, r *http.Request) {
	if "" == html {
		wd, _ := os.Getwd()
		tData := templateData{wd}
		var htmlBuf bytes.Buffer
		c := amber.New()
		c.Options.PrettyPrint = true
		c.Parse(amberTempl)

		t := template.Must(c.Compile())
		t.Delims("[{", "}]")
		err := t.Execute(&htmlBuf, tData)

		if nil != err {
			// log.Println(err)
		}

		html = htmlBuf.String()
	}
	fmt.Fprintf(w, html)
}

type User struct {
	Id    int
	Email string
	Name  string
}

var userStore = map[int]User{
	1: {Id: 1, Email: "tyrion@lannister.ws", Name: "Tyrion Lannister"},
}

var deletedUsers = []User{}

type UserService struct {
	gorest.RestService `root:"/service/" consumes:"application/json" produces:"application/json"`

	//End-Point level configs: Field names must be the same as the corresponding method names,
	// but not-exported (starts with lowercase)

	userDetails gorest.EndPoint `method:"GET" path:"/users/{Id:int}" output:"User"`
	listUsers   gorest.EndPoint `method:"GET" path:"/users/" output:"[]User"`
	addUser     gorest.EndPoint `method:"POST" path:"/users/" postdata:"User"`
	deleteUser  gorest.EndPoint `method:"DELETE" path:"/users/{Id:int}"`

	uptime		gorest.EndPoint `method:"GET" path:"/uptime" output:"string"`
}

func (srv UserService) Uptime() (uptime string) {
	return string(time.Now().Sub(startTime))
}

func (serv UserService) UserDetails(id int) (u User) {
	if user, found := userStore[id]; found {
		return user
	}
	serv.ResponseBuilder().SetResponseCode(404).Overide(true)
	// Overide causes the entity returned by the method to be ignored.
	// Other wise it would send back zeroed object
	return
}

func (serv UserService) ListUsers() []User {
	// serv.ResponseBuilder().CacheMaxAge(60 * 60 * 24)
	// List cacheable for a day. More work to come on this, Etag, etc
	users := make([]User, 0)
	for _, u := range userStore {
		users = append(users, u)
	}

	return users
}

func getNextUserId() (n int) {
	n = 1
	for _, u := range userStore {
		if u.Id > n {
			n = u.Id
		}
	}
	return n+1
}

func (serv UserService) AddUser(u User) {
	// log.Printf("userStore: %q\n", userStore)
	if u.Id == 0 {
		u.Id = getNextUserId() // Next id
	}

	// log.Printf("Adding user: %#q\n", u)

	userStore[u.Id] = u

	serv.ResponseBuilder().Created("http://" + address + serviceRoot + "users/" + string(u.Id))
	// Created, http 201
}

func (serv UserService) DeleteUser(id int) {
	if _, ok := userStore[id]; ok {
		delete(userStore, id)
		return
	}

	serv.ResponseBuilder().SetResponseCode(404).Overide(true)
	return
}
