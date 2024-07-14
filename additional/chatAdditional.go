package additional

import (
	"bytes"
	"chat/variables"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"time"
)

// func analyze auth status and return response view
func AuthRegView(c echo.Context, page string) error {
	var theme = ""
	toggleParam := c.QueryParam("toggle")
	if toggleParam == "true" {
		theme = Toggle(c)
	}
	if ValidateToken(c) == "Token error" {
		if theme == "" {
			theme = CheckTheme(c)
		}
		if theme == "dark" {
			return c.Render(http.StatusOK, page+"DARK.html", nil)
		} else {
			return c.Render(http.StatusOK, page+".html", nil)
		}

	} else {
		return c.Redirect(http.StatusMovedPermanently, "/chat")
	}
}

// func for chat views
func ParseMessages(c echo.Context, usernameID string, recipientID string) []variables.Message {
	_, err := variables.Db.Exec("insert into connects(userA, userB) values (?,?)", usernameID, recipientID)
	if err != nil {
		fmt.Println(err)
	}

	var messages []variables.Message
	result, err := variables.Db.Query("select messageID,text,file,fileName,cast(date as char) from message join content on content.contentID = message.contentID where message.fromID = ? and message.toID = ?", usernameID, recipientID)
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
		var elem variables.Message
		if bytes.Equal(byteFile, []byte{48}) == true {
			elem = variables.Message{Favourite: false, MessageID: id, Text: text, File: byteFile, FileName: "", Date: date, DateStr: date.Format("2006-01-02 15:04:05"), Pos: "r"}
		} else {
			elem = variables.Message{Favourite: false, MessageID: id, Text: text, File: byteFile, FileName: fileName, Date: date, DateStr: date.Format("2006-01-02 15:04:05"), Pos: "r"}
		}
		messages = append(messages, elem)
	}

	result, err = variables.Db.Query("select messageID,text,file,fileName,cast(date as char) from message join content on content.contentID = message.contentID where message.fromID = ? and message.toID = ?", recipientID, usernameID)
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
		var elem variables.Message
		if bytes.Equal(byteFile, []byte{48}) == true {
			elem = variables.Message{Favourite: false, MessageID: id, Text: text, File: byteFile, FileName: "", Date: date, DateStr: date.Format("2006-01-02 15:04:05"), Pos: "l"}
		} else {
			elem = variables.Message{Favourite: false, MessageID: id, Text: text, File: byteFile, FileName: fileName, Date: date, DateStr: date.Format("2006-01-02 15:04:05"), Pos: "l"}
		}
		messages = append(messages, elem)
	}
	sort.Slice(messages, func(i, j int) bool { return messages[i].Date.Before(messages[j].Date) })
	lastMessageCookie := &http.Cookie{}
	lastMessageCookie.Name = "lastMessage"
	if len(messages) > 0 {
		lastMessageCookie.Value = messages[len(messages)-1:][0].DateStr
	} else {
		now := time.Now()
		lastMessageCookie.Value = now.Format("2006-01-02 15:04:05")
	}
	lastMessageCookie.SameSite = 3
	lastMessageCookie.HttpOnly = true
	lastMessageCookie.Secure = true
	c.SetCookie(lastMessageCookie)
	return messages
}
func UpdateTimeVisit(userA string, userB string) {
	now := time.Now()
	lastVisitTime := now.Format("2006-01-02 15:04:05")
	_, err := variables.Db.Exec("update connects set timeOpen = ? where userA = ? and userB = ?",
		lastVisitTime, userA, userB)
	if err != nil {
		fmt.Println(err)
	}
}
func SearchRecipient(c echo.Context) string {
	var recipientID string
	recipient := c.QueryParam("recipient")
	recipientID = CheckUser(recipient)
	if recipientID != "" {
		newRecipientCookie := &http.Cookie{}
		newRecipientCookie.Name = "recipient"
		newRecipientCookie.Value = recipientID
		newRecipientCookie.SameSite = 3
		newRecipientCookie.HttpOnly = true
		newRecipientCookie.Secure = true
		c.SetCookie(newRecipientCookie)
	} else {
		var recipientCookie, err = c.Cookie("recipient")
		if err != nil {
			fmt.Println(err)
		}
		if recipientCookie != nil {
			recipientID = recipientCookie.Value
		}
	}
	return recipientID
}
func DeleteMessage(usernameID string, deleteID string) {
	_, err := variables.Db.Exec("delete from message where messageID = ? and "+
		"(message.fromID = ? or message.toID = ?)", deleteID, usernameID, usernameID)
	if err != nil {
		fmt.Println(err)
	}
}

func DeleteFavouriteMessage(usernameID string, deleteID string) {
	_, err := variables.Db.Exec("delete from favouritemessage where id = ? and favouriteuserid = ?",
		deleteID, usernameID)
	if err != nil {
		fmt.Println(err)
	}

}

func FavouriteMessage(usernameID string, favouriteID string) {
	result, err := variables.Db.Query("select fromID from message where(message.fromID = ? or message.toID = ?) and message.messageID = ?", usernameID, usernameID, favouriteID)
	if err != nil {
		fmt.Println(err)
	}
	var fromID string
	for result.Next() {
		err = result.Scan(&fromID)
	}
	_, err = variables.Db.Exec("insert into favouritemessage(fromID, messageID, favouriteuserid) values (?, ?, ?)", fromID, favouriteID, usernameID)
	if err != nil {
		fmt.Println(err)
	}
}
func DownloadFavourite(c echo.Context, favouriteID string) []variables.FavouriteMessages {
	var favouriteArray []variables.FavouriteMessages
	result, err := variables.Db.Query("select favouritemessage.id,username,favouritemessage.fromID,text,"+
		"cast(date as char),file,fileName from favouritemessage join message on "+
		"favouritemessage.messageID = message.messageID join chat.content on message.contentID = "+
		"content.contentID join user on favouritemessage.fromID = user.userID where favouritemessage.favouriteuserid "+
		"= ? and favouritemessage.favouriteuserid != favouritemessage.fromID", favouriteID)
	if err != nil {
		fmt.Println(err)
	}
	for result.Next() {
		var username, text, dateStr, fileName string
		var byteFile []byte
		var fromid, id int
		err = result.Scan(&id, &username, &fromid, &text, &dateStr, &byteFile, &fileName)
		date, err := time.Parse(variables.TimeFormat, dateStr)
		if err != nil {
			fmt.Println(err)
		}
		var elem variables.FavouriteMessages
		if bytes.Equal(byteFile, []byte{48}) == true {
			elem = variables.FavouriteMessages{MessageID: id, Favourite: true, UserName: username, FromID: fromid,
				Text: text, File: byteFile, FileName: "", Date: date,
				DateStr: date.Format("2006-01-02 15:04:05"), Pos: "l"}
		} else {
			elem = variables.FavouriteMessages{MessageID: id, Favourite: true, UserName: username, FromID: fromid,
				Text: text, File: byteFile, FileName: fileName, Date: date,
				DateStr: date.Format("2006-01-02 15:04:05"), Pos: "l"}
		}
		favouriteArray = append(favouriteArray, elem)
	}

	result, err = variables.Db.Query("select favouritemessage.id,username,favouritemessage.fromID,text"+
		",cast(date as char),file,fileName from favouritemessage join message "+
		"on favouritemessage.messageID = message.messageID join chat.content on "+
		"message.contentID = content.contentID join user on favouritemessage.fromID "+
		"= user.userID where favouritemessage.favouriteuserid = ? and favouritemessage.favouriteuserid "+
		"= favouritemessage.fromID", favouriteID)
	if err != nil {
		fmt.Println(err)
	}

	for result.Next() {
		var username, text, dateStr, fileName string
		var byteFile []byte
		var fromid, id int
		err = result.Scan(&id, &username, &fromid, &text, &dateStr, &byteFile, &fileName)
		if err != nil {
			fmt.Println(err)
		}
		date, err := time.Parse(variables.TimeFormat, dateStr)
		if err != nil {
			fmt.Println(err)
		}
		var elem variables.FavouriteMessages
		if bytes.Equal(byteFile, []byte{48}) == true {
			elem = variables.FavouriteMessages{MessageID: id, Favourite: true, UserName: username,
				FromID: fromid, Text: text, File: byteFile, FileName: "", Date: date,
				DateStr: date.Format("2006-01-02 15:04:05"), Pos: "r"}
		} else {
			elem = variables.FavouriteMessages{MessageID: id, Favourite: true, UserName: username,
				FromID: fromid, Text: text, File: byteFile, FileName: fileName, Date: date,
				DateStr: date.Format("2006-01-02 15:04:05"), Pos: "r"}
		}
		favouriteArray = append(favouriteArray, elem)
	}
	sort.Slice(favouriteArray, func(i, j int) bool { return favouriteArray[i].Date.Before(favouriteArray[j].Date) })
	return favouriteArray
}
