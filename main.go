package main

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", login)
	http.HandleFunc("/tester", tester)
	http.ListenAndServe(":8080", nil)
}

func login(w http.ResponseWriter, r *http.Request) {

	// application runs from local server
	// no security is built
	// fake login to prevent some employees from using this tool
	cookie, _ := r.Cookie("pses")
	if cookie != nil && cookie.Value == "4q55Rc3bo61" {
		http.Redirect(w, r, "/tester", 301)
		return
	}

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		params := strings.Split(string(body), "&")
		if len(params) != 1 {
			http.Error(w, "Incorrect request body", http.StatusInternalServerError)
		}
		println("Requesting: " + string(body))

		pass := strings.Split(params[0], "=")[1]
		if pass == "4J76sr59N1c" { // unencrypted simple login check | localhost
			expiration := time.Now().Add(2 * time.Hour)
			cookie := http.Cookie{Name: "pses", Value: "4q55Rc3bo61", Expires: expiration}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/tester", 301)
			return
		}

		w.Write([]byte("Login failed."))
		return
	}

	html := "<!DOCTYPE HTML><html><head><title>Stresser - Login</title></head><body>"
	html += "<form action=\"/\" method=\"post\">"
	html += "<table><tbody><tr><td>Password:</td>"
	html += "<td><input type=\"password\" name=\"passw\" value=\"\" /></td></tr>"
	html += "<tr><td>&nbsp;</td><td><input type=\"submit\" value=\"Login\" /><td></tr>"
	html += "</tbody></table></body></html>"

	w.Write([]byte(html))
}

func tester(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("pses")
	if cookie == nil || cookie.Value != "4q55Rc3bo61" {
		http.Redirect(w, r, "/", 301)
		return
	}

	plan := 1
	threads := 30000
	rampup := 120
	post := ""

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		params := strings.Split(string(body), "&")
		if len(params) != 3 {
			http.Error(w, "Incorrect request body", http.StatusInternalServerError)
		}
		println("Requesting: " + string(body))

		plan, _ = strconv.Atoi(strings.Split(params[0], "=")[1])
		threads, _ = strconv.Atoi(strings.Split(params[1], "=")[1])
		rampup, _ = strconv.Atoi(strings.Split(params[2], "=")[1])

		if plan != 1 && plan != 2 && plan != 3 && plan != 4 && plan != 5 {
			plan = 1
		}
		if threads < 15 && threads > 100000 {
			threads = 15
		}
		if rampup < 1 && rampup > 900 {
			rampup = 120
		}

		cl := "/home/ec2-user/stresser/"
		cl += "test.sh /home/ec2-user/streser/test_plan.jmx"
		cl += " " + strconv.Itoa(threads)
		cl += " " + strconv.Itoa(rampup)
		cl += " 1 &"

		script := "/home/ec2-user/stresser/"
		script += "test.sh"

		arg1 := "/home/ec2-user/stresser/test.jmx"
		arg2 := strconv.Itoa(threads)
		arg3 := strconv.Itoa(rampup)

		println("Trying to run: " + cl)

		if err := exec.Command(script, arg1, arg2, arg3, "1").Run(); err != nil {
			post = "<b>An error occured</b>"
			println(err.Error())
		} else {
			post = "<b>Test is sent!</b> This tool will not show you if the test is over or not. Please check AWS and Queue."
		}
	}

	html := "<!DOCTYPE HTML><html><head><title>Stresser</title></head><body>"
	html += "<span style=\"color:#ff0000\">Before you send a test, make sure nobody is testing</span><hr />"
	html += "<form action=\"/tester\" method=\"post\">"
	html += "<table><thead><tr><th>Plan</th><th>Users</th><th>Ramp-up</th><th>&nbsp;</th></tr></thead><tbody><tr>"
	html += "<td><select name=\"plan\"><option value=\"1\" selected=\"selected\" title=\"Test 1\">Test 1</option></select></td>"
	html += "<td><input type=\"text\" name=\"threads\" value=\"" + strconv.Itoa(threads) + "\" /></td>"
	html += "<td><input type=\"text\" name=\"rampup\" value=\"" + strconv.Itoa(120) + "\" /></td>"
	html += "<td><i>Hold mouse on an option to read the test description</i></td>"
	html += "</tr></tbody></table>"
	html += "<hr /><input type=\"submit\" value=\"Start\" onclick=\"document.getElementById('post').innerHTML='';\" />&nbsp;&nbsp;&nbsp;"
	html += "<span id=\"post\">" + post + "</span>"
	html += "</body></html>"

	w.Write([]byte(html))
}
