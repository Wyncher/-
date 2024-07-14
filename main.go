package main

import (
	"chat/additional"
	"chat/handlersPackage"
	"chat/logger"
	"chat/variables"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"html/template"
)

// Handlers

func main() {
	variables.Db, _ = sql.Open("mysql", "mysql:mysql@/chat")
	e := echo.New()
	logger.NewLogger()              // new
	e.Use(logger.LoggingMiddleware) // ne
	e.HTTPErrorHandler = additional.CustomHTTPErrorHandler
	//e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("web/*.html")),
	}
	e.Static("/", "/web/")
	// Routes
	e.GET("/login", handlersPackage.Auth)
	e.POST("/login", handlersPackage.AuthPOST)
	e.GET("/authMain", handlersPackage.AuthMain)
	e.GET("/registration", handlersPackage.RegGET)
	e.POST("/registration", handlersPackage.RegPOST)
	e.GET("/chat", handlersPackage.ChatGET)
	e.POST("/chat", handlersPackage.ChatPOST)
	e.GET("/", handlersPackage.MainPage)
	e.POST("/refresh", handlersPackage.GetMessages)
	e.POST("/loaduserconnects", handlersPackage.LoaderUser)
	e.POST("/loadnewmessage", handlersPackage.LoadNewMessageCounter)
	// Restricted group
	r := e.Group("/")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(variables.JwtCustomClaims)
		},
		SigningKey:  variables.Secret,
		TokenLookup: "cookie:JWTToken",
	}
	r.Use(echojwt.WithConfig(config))

	logger.Logger.LogInfo().Msg(e.Start(":1323").Error())
	//e.Logger.Fatal(e.StartAutoTLS(":443"))
	//e.Start(":1323")
	//if l, ok := e.Logger.(*log.Logger); ok {
	//	l.SetHeader("${time_rfc3339} ${level}")
	//}
}
