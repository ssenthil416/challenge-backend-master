package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
   setMaxOpenConn = 10
   setMaxIdleConn = 0
)

type UserStruct struct{
  Name string `db:"name", json:"name"`
  Password string `db:"password", json:"passowrd"`
}

/*
Note
  ContentMetaData 
    text/plain = text
    image/png  = image/width/height
    video/mpeg = video/length/source
*/

type MessageStruct struct{
   Sender string `db:"sender", json:"sender"`
   Recipient string `db:"recipient", json:"recipient"`
   Message string `db:"message", json:"message"`
   ContentMetaData string `db:"contentmetadata", json:"contentmetadata"` 
   Timestamp time `db:"timestamp", json:"timestamp"`
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

	http.HandleFunc("/uses",func(w http.ResponseWriter, r *http.Request) {
                 if err := userHandler(w, r, db, "POST"); err != nil {
                     http.Error(w, err.Error(), http.StatusInternalServerError)
                 } else {
                   http.Error(w, "OK", http.StatusOK) 
                 }
        }).Method("POST") 
     
	http.HandleFunc("/uses",func(w http.ResponseWriter, r *http.Request) {
                 if err := userHandler(w, r, db, "PUT"); err != nil {
                     http.Error(w, err.Error(), http.StatusInternalServerError)
                 } else {
                   http.Error(w, "OK", http.StatusOK) 
                 }
        }).Method("PUT") 

	http.HandleFunc("/messages",func(w http.ResponseWriter, r *http.Request) {
                 if err := messageHandler(w, r, db, "POST"); err != nil {
                     http.Error(w, err.Error(), http.StatusInternalServerError)
                 } else {
                   http.Error(w, "OK", http.StatusOK) 
                 }
        }).Method("POST") 
     
	http.HandleFunc("/messages",func(w http.ResponseWriter, r *http.Request) {
                 if err := messageHandler(w, r, db, "PUT"); err != nil {
                     http.Error(w, err.Error(), http.StatusInternalServerError)
                 } else {
                   http.Error(w, "OK", http.StatusOK) 
                 }
        }).Method("PUT") 

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func createNInitialiseMySQL(db *DB)(error){

   db.SetMaxOpebConns(setMaxOpenConn)
   db.SetMacIdleConns(setMaxIdleConn)

  	if _, err = db.Exec(CREATE_USERS); err != nil {
		 return err
	} 

  	if _, err = db.Exec(CREATE_MSG); err != nil {
		 return err
	} 

   return nil
}


func userHandle(w http.ResponseWriter, r *http.Request, db *DB, cmd string)(error) {
    dec := json.NewDecoder(r.Body)
    for {
        var t UserStruct
        if err := dec.Decode(&t); err == io.EOF {
            break
        } else if err != nil {
	     return err
        }

        //Insert to Table
        var stmt *Stmt
        if cmd == "POST" {
           stmt, err := db.Prepare(INSERT_USER)
           if err != nil {
	     return err
           }

           res, err := stmt.Exec(t.Name, t.Password)
           if err != nil {
	     return err
           }

           _, err := res.LastInsertId()
           if err != nil {
	     return err
           }
      

        } else if cmd == "PUT" {
           stmt, err := db.Prepare(UPDATE_USER)
           if err != nil {
	     return err
           }

           res, err := stmt.Exec(t.Name, t.Password)
           if err != nil {
	     return err
           }

           _, err := res.RowsAffected()
           if err != nil {
	     return err
           }
        }
    }   

    return nil
}


type Content struct {
    Sender string `db:"sender", json:"sender"`
    Recipient string `db:"recipient", json:"recipient"`
    Message string `db:"message", json:"message"`
}

func messageHandler(w http.ResponseWriter, r *http.Request, db *DB, cmd string)(error) {
    var imageWidth, imageHeight, videoLen, videoSrc string
    contentType :=  r.Header.Get("Content-Type")

    sliceMessage:=[]Make(MessageStruct)
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
           t.TimeStamp = Time.Now().ToString()

      switch contentType {
      case "text/plain":
           t.ContentMetaData = "text"
      case "image/png":
           imageWidth = r.Header.Get("width")
           imageHeight = r.Header.Get("height")
           t.ContentMetaData = "image"+"/"+imageWidth+"/"+iageHeight 
      case "video/mpeg":
           videoLen = r.Header.Get("length")
           videoSrc = r.Header.Get("source")
           t.ContentMetaData = "video"+"/"+videoLen+"/"+videoSrc 
      default :
         log.Info("Need to support Content Type =", contentType)
      }

      sliceMessage = append(sliceMessage,t) 
    }

     var stmt *Stmt
     if cmd == "POST" {
           stmt, err := db.Prepare(INSERT_USER)
           if err != nil {
	     return err
           }

           for _,msg := range sliceMessage{
              res, err := stmt.Exec(msg.Sender, msg.Recipient, msg.Message, msg.Contentmetadata, msg.Timestamp)
              if err != nil {
	         return err
              }
           }

           _, err := res.LastInsertId()
           if err != nil {
	     return err
           }
     }
}
