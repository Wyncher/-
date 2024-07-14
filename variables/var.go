package variables

import "database/sql"

// const variables
var TimeFormat = "2006-01-02 15:04:05"
var Db *sql.DB

// Secret phrase for encryption
var Secret = []byte("secret")
var AdditionalString = "aditionalS 97treng.6123"
