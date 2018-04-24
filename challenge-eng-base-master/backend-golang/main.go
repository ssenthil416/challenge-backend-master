package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
        "errors"

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
	Sender          string `db:"sender", json:"sender"`
	Recipient       string `db:"recipient", json:"recipient"`
	Message         string `db:"message", json:"message"`
	ContentMetaData string `db:"contentmetadata", json:"contentmetadata"`
	Timestamp       string `db:"timestamp", json:"timestamp"`
}

func main() {
	db, err := sql.Open("mysql", "root:testpass@tcp(db:3306)/challenge")
	if err != nil {
		log.Fatal("unable to connect to DB", err)
	}

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
                     fmt.Println("err =", err)
		     http.Error(w, "OK", http.StatusOK)
		   }
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
                 w.Header().Add("Content-Type", "application/json")
		if err := messageHandler(w, r, db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "OK", http.StatusOK)
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

func userHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) (err error) {
        cmd := r.Method
        firstTime := true
        fmt.Println("cmd =", cmd)
        fmt.Println("Content Type =", r.Header.Get("Conetent-Type"))
	dec := json.NewDecoder(r.Body)
        fmt.Printf("Body = %v\n", dec)
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

                fmt.Println("Data t =%v", t)
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

			res, err := stmt.Exec(t.Password,t.Name)
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
	Sender    string `db:"sender", json:"sender"`
	Recipient string `db:"recipient", json:"recipient"`
	Message   string `db:"message", json:"message"`
}

func messageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB)(err  error) {
	var imageWidth, imageHeight, videoLen, videoSrc string
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

		t.Message = content.Message
		t.Sender = content.Sender
		t.Recipient = content.Recipient
		t.Timestamp = time.Now().String()

		switch contentType {
		case "text/plain":
			t.ContentMetaData = "text"
		case "image/png":
			imageWidth = r.Header.Get("width")
			imageHeight = r.Header.Get("height")
			t.ContentMetaData = "image" + "/" + imageWidth + "/" + imageHeight
		case "video/mpeg":
			videoLen = r.Header.Get("length")
			videoSrc = r.Header.Get("source")
			t.ContentMetaData = "video" + "/" + videoLen + "/" + videoSrc
		default:
			log.Printf("Need to support Content Type =%v", contentType)
		}

		sliceMessage = append(sliceMessage, t)
	}

	var stmt *sql.Stmt
	if cmd == "POST" {
		stmt, err = db.Prepare(INSERT_USER)
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
