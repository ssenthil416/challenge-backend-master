package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	setMaxOpenConn = 10
	setMaxIdleConn = 0
)

type UserStruct struct {
	Name     string `db:"name", json:"name"`
	Password string `db:"password", json:"passowrd"`
}

/*
Note
  ContentMetaData
    text/plain = text
    image/png  = image/width/height
    video/mpeg = video/length/source
*/

type MessageStruct struct {
	Sender          string    `db:"sender", json:"sender"`
	Recipient       string    `db:"recipient", json:"recipient"`
	Message         string    `db:"message", json:"message"`
	ContentMetaData string    `db:"contentmetadata", json:"contentmetadata"`
	Timestamp       time.Time `db:"timestamp", json:"timestamp"`
}

func cleanUp(db *sql.DB) error {

	if _, err := db.Exec(DROP_MSG); err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := sql.Open("mysql", "root:testpass@tcp(db:3306)/challenge")
	if err != nil {
		log.Fatal("unable to connect to DB", err)
	}

	//DROP TABLE, if needed
	//cleanUp(db)
	//return

	if err := createNInitialiseMySQL(db); err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		var result string
		if err := db.QueryRow(`SELECT col FROM test`).Scan(&result); err != nil {
			log.Panic(err)
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"result":  result,
			"backend": "go",
		}); err != nil {
			log.Panic(err)
		}
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if err := userHandler(w, r, db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if err := messageHandler(w, r, db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func createNInitialiseMySQL(db *sql.DB) error {

	db.SetMaxOpenConns(setMaxOpenConn)
	db.SetMaxIdleConns(setMaxIdleConn)

	if _, err := db.Exec(CREATE_USERS); err != nil {
		return err
	}

	if _, err := db.Exec(CREATE_MSG); err != nil {
		return err
	}

	return nil
}

func userValidation(user UserStruct) bool {
	if user.Name == "" || user.Name == "null" || user.Name == "NULL" ||
		user.Password == "" || user.Password == "null" || user.Password == "NULL" {
		return false
	}

	return true
}

func userHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) (err error) {
	cmd := r.Method
	firstTime := true
	dec := json.NewDecoder(r.Body)
	for {
		var t UserStruct
		if err := dec.Decode(&t); err == io.EOF {
			if firstTime {
				return errors.New("Missing Params")
			}
			break
		} else if err != nil {
			return err
		}
		firstTime = false

		if userValidation(t) == false {
			return errors.New("Missing Params")
		}

		//Insert to Table
		var stmt *sql.Stmt
		if cmd == "POST" {
			stmt, err = db.Prepare(INSERT_USER)
			if err != nil {
				return err
			}

			res, err := stmt.Exec(t.Name, t.Password)
			if err != nil {
				return err
			}

			lid, err := res.LastInsertId()
			if err != nil {
				return err
			}

			fmt.Println("lastInsert ID =", lid)

		} else if cmd == "PUT" {
			stmt, err = db.Prepare(UPDATE_USER)
			if err != nil {
				return err
			}

			res, err := stmt.Exec(t.Password, t.Name)
			if err != nil {
				return err
			}

			reff, err := res.RowsAffected()
			if err != nil {
				return err
			}

			fmt.Println("RowsAffected =", reff)
		}
	}

	return nil
}

type Content struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
	Width     string `json:"width"`
	Height    string `json:"height"`
	Length    string `json:"length"`
	Source    string `json:"source"`
}

func msgValidation(c Content, msgType string) bool {
	if c.Sender == "" || c.Sender == "null" || c.Sender == "NULL" ||
		c.Recipient == "" || c.Recipient == "null" || c.Recipient == "NULL" ||
		c.Message == "" || c.Message == "null" || c.Message == "NULL" {
		return false
	}

	if msgType == "image" {
		if c.Width == "" || c.Width == "null" || c.Width == "NULL" ||
			c.Height == "" || c.Height == "null" || c.Height == "NULL" {
			return false
		}
	} else if msgType == "video" {
		if c.Length == "" || c.Length == "null" || c.Length == "NULL" ||
			c.Source == "" || c.Source == "null" || c.Source == "NULL" {
			return false
		}
	}
	return true
}

func messageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) (err error) {
	cmd := r.Method
	contentType := r.Header.Get("Content-Type")

	sliceMessage := make([]MessageStruct, 0, 1)
	dec := json.NewDecoder(r.Body)
	for {
		var t MessageStruct
		var content Content
		if err := dec.Decode(&content); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		//fmt.Println("contentType =", contentType)

		switch contentType {
		case "text/plain":
			if msgValidation(content, "text") == false {
				return errors.New("Missing Params")
			}
			t.ContentMetaData = "text"
		case "image/png":
			if msgValidation(content, "image") == false {
				return errors.New("Missing Params")
			}
			t.ContentMetaData = "image" + "/" + content.Width + "/" + content.Height
		case "video/mpeg":
			if msgValidation(content, "video") == false {
				return errors.New("Missing Params")
			}
			t.ContentMetaData = "video" + "/" + content.Length + "/" + content.Source
		default:
			log.Printf("Need to support Content Type =%v", contentType)
			return errors.New("Content Type Not Supported Yet")
		}

		t.Message = content.Message
		t.Sender = content.Sender
		t.Recipient = content.Recipient
		t.Timestamp = time.Now()

		sliceMessage = append(sliceMessage, t)
	}

	var stmt *sql.Stmt
	if cmd == "POST" {
		stmt, err = db.Prepare(INSERT_MSG)
		if err != nil {
			return err
		}
		var res sql.Result
		for _, msg := range sliceMessage {
			res, err = stmt.Exec(msg.Sender, msg.Recipient, msg.Message, msg.ContentMetaData, msg.Timestamp)
			if err != nil {
				return err
			}
		}

		lid, err := res.LastInsertId()
		if err != nil {
			return err
		}

		fmt.Println("LastInsertID =", lid)
	}

	return nil
}
