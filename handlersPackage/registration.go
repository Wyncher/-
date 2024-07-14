package handlersPackage

import (
	"bytes"
	"chat/additional"
	"chat/variables"
	"crypto/md5"
	"github.com/disintegration/imaging"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"image/jpeg"
	"net/http"
	"time"
)

func RegGET(c echo.Context) error {
	return additional.AuthRegView(c, "reg")
}

func RegPOST(c echo.Context) error {
	username := c.FormValue("username")
	result, err := variables.Db.Query("select username from user where username = ?", username)
	if err != nil {
		return err
	}
	if result.Next() != false {
		return err
	}
	email := c.FormValue("email")
	resultEmail, err := variables.Db.Query("select email from user where email = ?", email)
	if err != nil {
		return err
	}
	if resultEmail.Next() != false {
		return err
	}
	password := c.FormValue("password")
	logo, err := c.FormFile("logo")
	if err != nil {
		return err
	}
	src, err := logo.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	var dst []byte
	img, _ := jpeg.Decode(src)

	compressedImage := imaging.Resize(img, 100, 100, imaging.Lanczos)
	buf := new(bytes.Buffer)
	if err != nil {
		return err
	}
	err = jpeg.Encode(buf, compressedImage, nil)
	dst = buf.Bytes()
	if err != nil {
		return err
	}
	hashPassword := md5.Sum([]byte(password + variables.AdditionalString))
	_, err = variables.Db.Exec("insert into user (username, email, password,logo) values (?, ?, ?, ?)",
		username, email, string(hashPassword[:]), dst)
	if err != nil {
		return err
	}

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

	return c.Redirect(http.StatusMovedPermanently, "/chat")
}
