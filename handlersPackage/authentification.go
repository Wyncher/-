package handlersPackage

import (
	"chat/additional"
	"chat/variables"
	"crypto/md5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func Auth(c echo.Context) error {
	return additional.AuthRegView(c, "auth")
}
func AuthPOST(c echo.Context) error {
	if additional.ValidateToken(c) == "Token error" {
		username := c.FormValue("username")
		password := c.FormValue("password")

		hash := md5.Sum([]byte(password + variables.AdditionalString))
		result, err := variables.Db.Query("select password from user where username = ?", username)
		if err != nil {
			return err
		}

		for result.Next() {
			var u []byte
			err = result.Scan(&u)
			var fixedSizePassword [16]byte
			copy(fixedSizePassword[:], u)

			// Throws unauthorized error
			if fixedSizePassword == hash {

				// Set custom claims
				claims := &variables.JwtCustomClaims{
					Name: username,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
					},
				}

				// Create token with claims
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				// Generate encoded token and send it as response.
				t, err := token.SignedString(variables.Secret)
				if err != nil {
					return err
				}
				JWTCookie := &http.Cookie{}
				JWTCookie.Name = "JWTToken"
				JWTCookie.Expires = time.Now().Add(time.Hour * 72)
				JWTCookie.Value = t
				JWTCookie.SameSite = 3
				JWTCookie.HttpOnly = true
				JWTCookie.Secure = true
				c.SetCookie(JWTCookie)

			}
		}
	}
	return c.Redirect(http.StatusMovedPermanently, "/chat")

}
