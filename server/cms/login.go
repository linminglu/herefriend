package cms

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

// CmsDashBoard .
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
