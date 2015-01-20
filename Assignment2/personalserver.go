/*Copyright Joseph Sturtevant 1/11/15
Joseph Sturtevant
CSS 490 Tactical Programming
Assignment 1

This is a simple time server

Code modified from the wiki:
https://golang.org/doc/articles/wiki/

Advice on custom 404 from StackOverFlow user Mostafa
http://stackoverflow.com/questions/9996767/showing-custom-404-error-page-with-standard-http-package
*/
package main

import (
	"fmt"
    "net/http"
    "time"
    "flag"
    "bytes"
    "os/exec"
    "strings"
)

//Flag variables for port and version, as well as the current version
var (
	portFlag = flag.Int("port", 8080, "Defines the port number to listen on")
	versionFlag = flag.Bool("V", false, "Returns the version")
	version = "2.3"
	users = map[string]string{}
)

//Handles calls to /time/
//Formats the time to HH:MM:SS AM/PM
func timeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/time/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	t := time.Now()
	curTime := t.Format("3:04:05 PM")
    fmt.Fprintf(w, "<html><head><title>The Time</title></head>")
    fmt.Fprintf(w, "<body><p>The time is : <span style='color:red;font-size:2em'>%s</span></p></body></html>", curTime)
}

//Handles calls to pretty much everywhere other than /time
func homeHandler(w http.ResponseWriter, r *http.Request){
	//http.Redirect(w, r, "./login", 302)
	fmt.Printf("HOME HANDLER\n")
	if r.URL.Path != "/" && r.URL.Path != "/index.html" && r.URL.Path != "/index.htm" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	uid, err := r.Cookie("userid")
	if err != nil {
		fmt.Fprintf(w, "<html><body><head><title>Login</title></head>")
		fmt.Fprintf(w, "<body><form action='login'>What is your name, Earthling?<input type='text' name='name' size='50'>")
		fmt.Fprintf(w, "<input type='submit'></form><p/></body></html>")
	} else {
		s := strings.TrimPrefix(uid.String(), "userid=")
		fmt.Printf("%s\n", s)
		fmt.Fprintf(w, "<html><body><head><title>Logged In</title></head>")
		fmt.Fprintf(w, "<body>Greetings, %s</body></html>", users[s])
	}	
}

func logoutHandler(w http.ResponseWriter, r *http.Request){
	fmt.Printf("LOGOUT HANDLER\n")
	if r.URL.Path != "/logout/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	if _, err := r.Cookie("userid"); err != nil {
		fmt.Printf("The error was: %s", err)
	} else {
		http.SetCookie(w, &(http.Cookie{Name: "userid", Path: "/", Expires: time.Now()}))
	}
	fmt.Fprintf(w, "<html><body><head><title>Logout</title><META http-equiv='refresh' content='10;URL=/''></head>")
	fmt.Fprintf(w, "<body><p>Good-bye.</p></body></html>")
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	fmt.Printf("LOGIN HANDLER\n")
	name := r.FormValue("name")
	fmt.Printf("The name given was: %s\n", name)
	cmd := exec.Command("uuidgen")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	//This lines gets rid of the /n in out.String()
	s := strings.TrimSuffix(out.String(), "\n")
	users[s] = name
	fmt.Printf("User Name %s was stored in users[%s]\n", users[s], s)
	c := http.Cookie{Name: "userid", Value: out.String(), Path: "/"}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "./..", 302)
}

//Error handler
//Prints a custom page on StatusNotFound error (404)
func errorHandler(w http.ResponseWriter, r *http.Request, status int){
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprintf(w, "<html><head><title>You Dun Goofed</title></head>")
		fmt.Fprintf(w, "<body><p>These are not the URLs you're looking for.</p></body></html>")
	}
}

//Starts the server. 
//Doesn't run if the -V flag is set
func runServer(){
	fmt.Printf("SERVER STARTED\n")
	http.HandleFunc("/time/", timeHandler)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/logout/", logoutHandler)
	http.HandleFunc("/login", loginHandler)
    if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), nil); err != nil{
    	fmt.Printf("Port %v already in use", *portFlag)
    }
}

func main() {
	fmt.Printf("PROGRAM STARTED\n")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("Version: %v\n", version)
	} else {
		runServer()
	}
}