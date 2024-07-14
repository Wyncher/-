package handlersPackage

import (
	"chat/additional"
	"chat/variables"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
)

func LogoutFunc(c echo.Context) error {
	Cookie := &http.Cookie{}
	Cookie.Name = "JWTToken"
	Cookie.Value = ""
	Cookie.SameSite = 3
	Cookie.HttpOnly = true
	Cookie.Secure = true
	c.SetCookie(Cookie)
	Cookie.Name = "recipient"
	c.SetCookie(Cookie)
	Cookie.Name = "lastMessage"
	c.SetCookie(Cookie)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func FavouritePage(c echo.Context) {
	recipientCookie := &http.Cookie{}
	recipientCookie.Name = "recipient"
	recipientCookie.Value = ""
	recipientCookie.SameSite = 3
	recipientCookie.HttpOnly = true
	recipientCookie.Secure = true
	c.SetCookie(recipientCookie)
}

func ChatGET(c echo.Context) error {
	logoutParam := c.QueryParam("logout")
	if logoutParam != "" {
		return LogoutFunc(c)
	}
	var theme = ""
	toggleParam := c.QueryParam("toggle")
	if toggleParam == "true" {
		theme = additional.Toggle(c)
	}
	if theme == "" {
		theme = additional.CheckTheme(c)
	}
	username := additional.ValidateToken(c)
	if username != "Token error" {
		usernameID := additional.CheckUser(username)
		if usernameID == "" {
			return c.Redirect(http.StatusMovedPermanently, "/")
		}

		//search recipient MessageID from Cookie or GET request
		recipientID := additional.SearchRecipient(c)
		//if user press on favourite messages
		favouriteChatsParam := c.QueryParam("favouritechats")
		if favouriteChatsParam != "" {
			FavouritePage(c)
			recipientID = ""
		}
		//if user press on delete message button
		deleteFavouriteID := c.QueryParam("deletefavourite")
		if deleteFavouriteID != "" {
			additional.DeleteFavouriteMessage(usernameID, deleteFavouriteID)
		}
		deleteID := c.QueryParam("delete")
		if deleteID != "" {
			additional.DeleteMessage(usernameID, deleteID)
		}
		//if user press on favourite message button
		favouriteID := c.QueryParam("favourite")
		if favouriteID != "" {
			additional.FavouriteMessage(usernameID, favouriteID)
		}

		if recipientID != "" && recipientID != usernameID {
			additional.UpdateTimeVisit(usernameID, recipientID)
			messages := additional.ParseMessages(c, usernameID, recipientID)
			if theme == "dark" {
				return c.Render(http.StatusOK, "chatDARK.html", messages)
			} else {
				return c.Render(http.StatusOK, "chat.html", messages)
			}
		} else {
			favouriteMessagesArray := additional.DownloadFavourite(c, usernameID)
			if theme == "dark" {
				return c.Render(http.StatusOK, "chatDARK.html", favouriteMessagesArray)
			} else {
				return c.Render(http.StatusOK, "chat.html", favouriteMessagesArray)
			}
		}

	} else {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}

}

func ChatPOST(c echo.Context) error {
	username := additional.ValidateToken(c)
	usernameID := additional.CheckUser(username)
	if usernameID == "" {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	if c.FormValue("search") != "" {
		return LoaderUser(c)
	}
	recipientID := additional.SearchRecipient(c)
	//read text from form
	text := c.FormValue("message")
	// read file
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("no file")
	}
	var dst = []byte{48}
	var fName string
	if file == nil {
		//if file not appended to message
		fName = ""
	} else {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, _ = io.ReadAll(src)
		fName = "Вложение: " + file.Filename
	}
	//insert content of new message to DB
	_, err = variables.Db.Exec("insert into chat.content (text, file, fileName) values (?,?,?)", text, dst, fName)
	if err != nil {
		log.Error(err)
	}
	//insert new message to DB
	_, err = variables.Db.Exec("insert into message (fromID, toID, contentID) values (?,?,LAST_INSERT_ID());", usernameID, recipientID)
	if err != nil {
		log.Error(err)
	}
	return c.Redirect(http.StatusMovedPermanently, "/chat")
}
