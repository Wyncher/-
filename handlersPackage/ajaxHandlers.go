package handlersPackage

import (
	"bytes"
	"chat/additional"
	"chat/variables"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"sort"
	"time"
)

func LoadNewMessageCounter(c echo.Context) error {
	var newDate time.Time
	var messages []variables.MessagesCounter
	username := additional.ValidateToken(c)
	if username != "Token error" {
		usernameID := additional.CheckUser(username)
		result, err := variables.Db.Query("select id,userA,userB,timeOpen from connects where userA = ?", usernameID)
		if err != nil {
			fmt.Println(err)
		}
		for result.Next() {
			var dateStr string
			var id, userA, userB int
			err = result.Scan(&id, &userA, &userB, &dateStr)
			date, err := time.Parse(variables.TimeFormat, dateStr)
			if err != nil {
				fmt.Println(err)
			}
			resMessage, err := variables.Db.Query("select cast(date as char) from message where message.fromID = ? and message.toID = ? and date > ?", userB, userA, newDate.Format("2006-01-02 15:04:05"))
			var counter = 0
			for resMessage.Next() {
				var strDateMessage string
				err = resMessage.Scan(&strDateMessage)
				dateMessage, _ := time.Parse(variables.TimeFormat, strDateMessage)
				if dateMessage.After(date) {
					counter++
					date = dateMessage
				}

			}
			if counter != 0 {
				var elem variables.MessagesCounter
				elem = variables.MessagesCounter{UserID: userB, MCounter: counter}
				messages = append(messages, elem)

			}

		}

	}

	return c.JSON(http.StatusOK, messages)
}

func LoaderUser(c echo.Context) error {
	var result *sql.Rows
	var err error
	username := additional.ValidateToken(c)
	usernameID := additional.CheckUser(username)
	if usernameID == "" {
		return c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		result, err = variables.Db.Query("select user.UserID,username,logo from connects join user on connects.userB = user.userID and connects.userA = ? limit 10", usernameID)
		if err != nil {
			fmt.Println(err)
		}
	}

	if c.FormValue("search") != "" {
		// retrieve search string
		queryString := c.FormValue("search")

		//ВНИМАЕНИЕ

		result, err = variables.Db.Query("select user.UserID,username,logo from user where username like '%" + queryString + "%' limit 10 ")

		//ВНИМАНИЕ
	}
	if err != nil {
		fmt.Println(err)
	}

	var users []variables.User
	for result.Next() {
		var uid int
		var username string
		var logo []byte
		err = result.Scan(&uid, &username, &logo)
		strLogo := base64.StdEncoding.EncodeToString(logo)

		users = append(users, variables.User{Id: uid, Username: username, Logo: strLogo})
	}
	if err != nil {
		log.Error(err)
	}
	return c.JSON(http.StatusOK, users)
}

func GetMessages(c echo.Context) error {
	lastMessage, err := c.Cookie("lastMessage")
	var newDate time.Time
	var newDateForCookie time.Time
	if err != nil {
		newDateForCookie, err = time.Parse(variables.TimeFormat, "2000-01-01 00:00:00")
		if err != nil {
			fmt.Println(err)
		}
	} else {
		if lastMessage.Value != "" {
			newDateForCookie, err = time.Parse(variables.TimeFormat, lastMessage.Value)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	var messages []variables.Message
	newDate = newDateForCookie
	username := additional.ValidateToken(c)
	if username != "Token error" {
		recipientID := additional.SearchRecipient(c)
		usernameID := additional.CheckUser(username)
		result, err := variables.Db.Query("select messageID,text,file,fileName,cast(date as char) from message "+
			"join content on content.contentID = message.contentID join user on user.userID = message.fromID "+
			"where message.fromID = ? and message.toID = ? and date > ?", usernameID, recipientID,
			newDate.Format("2006-01-02 15:04:05"))
		if err != nil {
			fmt.Println(err)
		}
		for result.Next() {
			var text, dateStr, fileName string
			var byteFile []byte
			var id int
			err = result.Scan(&id, &text, &byteFile, &fileName, &dateStr)
			date, err := time.Parse(variables.TimeFormat, dateStr)
			if err != nil {
				fmt.Println(err)
			}
			if newDate.After(date) {

			} else {
				newDateForCookie = date
				var elem variables.Message
				if bytes.Equal(byteFile, []byte{48}) == true {
					elem = variables.Message{Favourite: false, MessageID: id, Text: text, File: byteFile, FileName: "", Date: date, DateStr: date.Format("2006-01-02 15:04:05"), Pos: "r"}
				} else {
					elem = variables.Message{Favourite: false, MessageID: id, Text: text, File: byteFile, FileName: fileName, Date: date, DateStr: date.Format("2006-01-02 15:04:05"), Pos: "r"}
				}
				messages = append(messages, elem)
			}
		}
		result = nil
		result, err = variables.Db.Query("select messageID,text,file,fileName,cast(date as char) from message join content on content.contentID = message.contentID join user on user.userID = message.fromID where message.fromID = ? and message.toID = ? and date > ?", recipientID, usernameID, newDate.Format("2006-01-02 15:04:05"))
		if err != nil {
			fmt.Println(err)
		}

		for result.Next() {
			var text, dateStr, fileName string
			var byteFile []byte
			var id int
			err = result.Scan(&id, &text, &byteFile, &fileName, &dateStr)
			date, err := time.Parse(variables.TimeFormat, dateStr)
			if err != nil {
				fmt.Println(err)
			}
			if newDate.After(date) {

			} else {
				newDateForCookie = date
				var elem variables.Message
				if bytes.Equal(byteFile, []byte{48}) == true {
					elem = variables.Message{Favourite: false, MessageID: id, Text: text, File: byteFile, FileName: "", Date: date, DateStr: date.Format("2006-01-02 15:04:05"), Pos: "l"}
				} else {
					elem = variables.Message{Favourite: false, MessageID: id, Text: text, File: byteFile, FileName: fileName, Date: date, DateStr: date.Format("2006-01-02 15:04:05"), Pos: "l"}
				}
				messages = append(messages, elem)
			}
		}
		sort.Slice(messages, func(i, j int) bool { return messages[i].Date.Before(messages[j].Date) })
		lastMessageCookie := &http.Cookie{}
		lastMessageCookie.Name = "lastMessage"
		lastMessageCookie.Value = newDateForCookie.Format("2006-01-02 15:04:05")
		lastMessageCookie.SameSite = 3
		lastMessageCookie.HttpOnly = true
		lastMessageCookie.Secure = true
		c.SetCookie(lastMessageCookie)
	}

	return c.JSON(http.StatusOK, messages)

}
