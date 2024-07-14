package additional

import (
	"chat/variables"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// fonc for validate JWT token
func ValidateToken(c echo.Context) string {
	//retrieve cookie value
	cookie, _ := c.Cookie("JWTToken")
	if cookie == nil {
		return "Token error"
	}
	//try parse token
	token, _ := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return variables.Secret, nil
	})
	//if token not valid
	if token == nil {
		return "Token error"
	}
	//try parse claim from token
	claims, ok := token.Claims.(jwt.MapClaims)
	//if token claims not valid
	if !ok {
		return "Token error"
	}
	//return username string from claim
	return claims["name"].(string)
}

// func retrieve Theme from cookie
func CheckTheme(c echo.Context) string {
	//retrieve cookie value
	theme, err := c.Cookie("theme")
	if err != nil {
		//if cookie not found
		return "ERR"
	}
	//return theme cookie value
	return theme.Value
}

// func check if exist username in DB
func CheckUser(username string) string {
	var ID string
	//query for DB
	result, _ := variables.Db.Query("select userID from user where username = ? or userID = ?", username, username)
	if result == nil {
		//if username not found
		return ""
	}
	for result.Next() {
		//scan userID
		err := result.Scan(&ID)
		if err != nil {
			fmt.Println(err)
		}
	}
	//return username MessageID
	return ID
}

// function for Toggle theme value on cookie
func Toggle(c echo.Context) string {
	//make cookie
	theme, err := c.Cookie("theme")
	Cookie := &http.Cookie{}
	Cookie.Name = "theme"
	Cookie.SameSite = 3
	Cookie.HttpOnly = true
	Cookie.Secure = true
	//if theme cookie not found
	if err != nil {
		Cookie.Value = "light"
		c.SetCookie(Cookie)
		return Cookie.Value
	}
	//if last theme is dark
	if theme.Value == "dark" {
		Cookie.Value = "light"

	}
	//if last theme is light
	if theme.Value == "light" {
		Cookie.Value = "dark"

	}
	//set cookie
	c.SetCookie(Cookie)
	//return current theme
	return Cookie.Value
}

// handler for HTTP errors
func CustomHTTPErrorHandler(err error, c echo.Context) {
	//retrieve code error
	code := http.StatusInternalServerError
	//if he, ok := err.(*echo.HTTPError); ok {
	//	code = he.Code
	//}
	//c.Logger().Error(err)
	theme := CheckTheme(c)
	errorPage := ""
	//make error page lin
	if theme == "dark" {
		errorPage = fmt.Sprintf("web/errors/%dDARK.html", code)
	} else {
		errorPage = fmt.Sprintf("web/errors/%d.html", code)
	}

	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
}
