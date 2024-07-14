package handlersPackage

import (
	"chat/additional"
	"github.com/labstack/echo/v4"
	"net/http"
)

func MainPage(c echo.Context) error {
	//make theme cookie
	Cookie := &http.Cookie{}
	Cookie.Name = "theme"
	Cookie.Value = "light" // light theme default
	Cookie.SameSite = 3
	Cookie.HttpOnly = true
	Cookie.Secure = true
	c.SetCookie(Cookie)
	//if authorize
	if additional.ValidateToken(c) == "Token error" {
		return c.Redirect(http.StatusMovedPermanently, "/authMain")
	} else {
		return c.Redirect(http.StatusMovedPermanently, "/chat")
	}
}
func AuthMain(c echo.Context) error {
	return additional.AuthRegView(c, "authMain")
}
