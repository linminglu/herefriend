package cms

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

/*
 |    Function: dashboard
 |      Author: Mr.Sancho
 |        Date: 2016-01-30
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func CmsDashBoard(c *gin.Context) {
	username := c.DefaultPostForm("username", "")
	if username != "sc" && username != "fan" {
		t, _ := template.ParseFiles("public/err500.html")
		t.Execute(c.Writer, nil)
		return
	}

	t, _ := template.ParseFiles("public/dashboard.html")
	t.Execute(c.Writer, nil)
	return
}
