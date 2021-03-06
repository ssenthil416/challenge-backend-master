package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	setMaxOpenConn = 10
	setMaxIdleConn = 0
	timeFormat     = "2006-01-02 15:04:05"
	userCacheSize  = 10
)

//UserStruct data for each users
type UserStruct struct {
	Name     string `db:"name", json:"name"`
	Password string `db:"password", json:"passowrd"`
}

/*
Note
  Contentmetadata
    text/plain = text
    image/png  = image/width/height
    video/mpeg = video/length/source
*/
//MessageStruct meta for Chat Messages
type MessageStruct struct {
	Sender          string `db:"sender", json:"sender"`
	Recipient       string `db:"recipient", json:"recipient"`
	Message         string `db:"message", json:"message"`
	Contentmetadata string `db:"contentmetadata", json:"contentmetadata"`
	Timestamp       string `db:"timestamp", json:"timestamp"`
}

//Cache
var (
	userCache = make([]string, userCacheSize)
)

//insertUserCache use to keep recent userCache Size in cache
func insertUserCache(usr string) {

	if len(userCache) == userCacheSize {
		tmp := userCache[1:]
		userCache = tmp
	}
	userCache = append(userCache, usr)
}

//CleanUp for Testing purpose
func cleanUp(db *sql.DB) error {

	if _, err := db.Exec(DROP_USER); err != nil {
		return err
	}

	if _, err := db.Exec(DROP_MSG); err != nil {
		return err
	}

	return nil
}

//Main Routine
func main() {
	db, err := sql.Open("mysql", "root:testpass@tcp(db:3306)/challenge")
	if err != nil {
		log.Fatal("unable to connect to DB", err)
	}

	defer db.Close()

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

	//Users rest end point
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if err := userHandler(w, r, db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	//MEssages rest end point
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if err := messageHandler(w, r, db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		/*else {
			w.WriteHeader(http.StatusOK)
		}*/
	})

	//webserver
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

//Create tables and set max conn to DB limits
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

//Simple user validation
func userValidation(user UserStruct) bool {
	if user.Name == "" || user.Name == "null" || user.Name == "NULL" ||
		user.Password == "" || user.Password == "null" || user.Password == "NULL" {
		return false
	}

	return true
}

//User Handler
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

		//Insert/Update to Table
		var stmt *sql.Stmt
		if cmd == "POST" {
			//insert to userCache
			insertUserCache(t.Name)

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

//Content to read message details per row
type Content struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
	Width     string `json:"width"`
	Height    string `json:"height"`
	Length    string `json:"length"`
	Source    string `json:"source"`
}

//Message simple validation
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

//MessabeHandler
func messageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) (err error) {
	cmd := r.Method
	contentType := r.Header.Get("Content-Type")

	if cmd == "GET" {
		err = msgGetHandler(w, r, db)
		return err
	}

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
			t.Contentmetadata = "text"
		case "image/png":
			if msgValidation(content, "image") == false {
				return errors.New("Missing Params")
			}
			t.Contentmetadata = "image" + "/" + content.Width + "/" + content.Height
		case "video/mpeg":
			if msgValidation(content, "video") == false {
				return errors.New("Missing Params")
			}
			t.Contentmetadata = "video" + "/" + content.Length + "/" + content.Source
		default:
			log.Printf("Need to support Content Type =%v", contentType)
			return errors.New("Content Type Not Supported Yet")
		}

		t.Message = content.Message
		t.Sender = content.Sender
		t.Recipient = content.Recipient
		tmpTime := time.Now()
		tmpTime.Format(timeFormat)
		t.Timestamp = tmpTime.String()

		sliceMessage = append(sliceMessage, t)
	}

	//Insert/Update message
	var stmt *sql.Stmt
	if cmd == "POST" {
		stmt, err = db.Prepare(INSERT_MSG)
		if err != nil {
			return err
		}
		var res sql.Result
		for _, msg := range sliceMessage {
			res, err = stmt.Exec(msg.Sender, msg.Recipient, msg.Message, msg.Contentmetadata, msg.Timestamp)
			if err != nil {
				return err
			}
		}

		lid, err := res.LastInsertId()
		if err != nil {
			return err
		}

		fmt.Println("LastInsertID =", lid)
	} else if cmd == "PUT" {
		stmt, err = db.Prepare(UPDATE_MSG)
		if err != nil {
			return err
		}
		var res sql.Result
		for _, msg := range sliceMessage {
			res, err = stmt.Exec(msg.Message, msg.Sender, msg.Recipient, msg.Timestamp)
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return nil
}

//Validate User
func validateUser(sender string, recipient string, db *sql.DB) (err error) {
	found := false
	s := false
	r := false

	//Check it in cache first
	for _, user := range userCache {
		if user == sender {
			s = true
		} else if user == recipient {
			r = true
		} else if s == true && r == true {
			found = true
			break
		}
	}
	if found {
		return nil
	}

	//Check users in data base
	scount := 0
	err = db.QueryRow("SELECT count(*) from users WHERE name=?", sender).Scan(&scount)
	if err != nil {
		return err
	}
	rcount := 0
	err = db.QueryRow("SELECT count(*) from users WHERE name=?", recipient).Scan(&rcount)
	if err != nil {
		return err
	}

        //fmt.Println(scount,rcount)
	if scount == 0 || rcount == 0 {
		return errors.New("Invalid Users")
	}

	return nil
}

// Get command handler for Message end point
func msgGetHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) (err error) {

	sender := r.URL.Query()["sender"][0]
	recipient := r.URL.Query()["recipient"][0]
	pagenumberstr := r.URL.Query()["pagenumber"][0]
	totalmsginpagestr := r.URL.Query()["totalmsginpage"][0]

	//Basic Validation
	if sender == "" || sender == "null" || sender == "NULL" ||
		recipient == "" || recipient == "null" || recipient == "NULL" ||
		pagenumberstr == "" || pagenumberstr == "null" || pagenumberstr == "NULL" ||
		totalmsginpagestr == "" || totalmsginpagestr == "null" || totalmsginpagestr == "NULL" {
		return errors.New("Missing Params")
	}
	pagenumber, _ := strconv.Atoi(pagenumberstr)
	totalmsginpage, _ := strconv.Atoi(totalmsginpagestr)

	//Validate sender and recipient
	if err := validateUser(sender, recipient, db); err != nil {
		return err
	}

	//fmt.Println(pagenumber, totalmsginpage)

	//Query and reads all messages
	rows, err := db.Query(SELECT_MSG, sender, recipient)
	if err != nil {
		fmt.Println("err =", err)
		return err
	}
	scanData := make([]MessageStruct, 0, 1)
	for rows.Next() {
		var t MessageStruct
		if err = rows.Scan(&t.Sender, &t.Recipient, &t.Message, &t.Contentmetadata, &t.Timestamp); err != nil {
			fmt.Println("err =", err)
			break
		}
		scanData = append(scanData, t)
	}

	totalData := len(scanData)
	if totalData == 0 {
		return errors.New("No Message Found")
	}

	//Calculate which part to show as per pagenumber and number of message in a page
	sIndex := 0
	reqData := totalmsginpage * pagenumber
	//fmt.Println("reqData =", reqData)
	startPoint := totalmsginpage * (pagenumber - 1)
	//fmt.Println("startPoint =", startPoint)
	//fmt.Println("totalData =", totalData)
	if totalData >= reqData {
		sIndex = totalData - startPoint
	} else {
		sIndex = totalData
	}

	//fmt.Println("sIndex =", sIndex)

	//Final Output data
	outData := make([]MessageStruct, 0, 1)

	for i := (sIndex - 1); i > (sIndex - totalmsginpage - 1); i-- {
		outData = append(outData, scanData[i])
	}
	w.Header().Add("Content-Type", "application/json")
	//fmt.Println("outData =", len(outData))
	if len(outData) > 0 {
		json.NewEncoder(w).Encode(outData)
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
