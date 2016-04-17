package page

import (
	"fmt"
	"net/http"

	"herefriend/lib"
)

func getAuthCookies() []*http.Cookie {
	var cookies []*http.Cookie

	resp, err := lib.Get("http://my.baihe.com/Getinterlogin/gotoLogin?jsonCallBack=jQuery1830007877889787778258_1452424575779&event=3&spmp=4.20.87.225.1049&txtLoginEMail=scisgood%40foxmail.com&txtLoginPwd=Sancho87&chkRememberMe=0", nil)
	if nil == err {
		defer resp.Body.Close()
		for _, c := range resp.Cookies() {
			if "AuthCookie" == c.Name {
				fmt.Printf("【Get AuthCookie】%s\r\n", c.Value)
				cookies = []*http.Cookie{
					&http.Cookie{Name: "AuthCookie", Value: c.Value, Path: "/"},
				}
			}
		}
	}

	return cookies
}
