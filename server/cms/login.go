package cms

import (
	"html/template"
	"net/http"
)

/*
 |    Function: login
 |      Author: Mr.Sancho
 |        Date: 2016-01-30
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func CmsLogin(w http.ResponseWriter) {
	t, _ := template.ParseFiles("public/signin.html")
	t.Execute(w, nil)
}

/*
 |    Function: dashboard
 |      Author: Mr.Sancho
 |        Date: 2016-01-30
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func CmsDashBoard(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username, ok := r.Form["username"]
	if true != ok || ("sc" != username[0] && "fan" != username[0]) {
		t, _ := template.ParseFiles("public/err500.html")
		t.Execute(w, nil)
		return
	}

	t, _ := template.ParseFiles("public/dashboard.html")
	t.Execute(w, nil)
}
