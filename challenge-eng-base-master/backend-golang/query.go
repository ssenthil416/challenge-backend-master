package main

const (

	CREATE_USER = `
        	CREATE TABLE IF NOT EXISTS users (
       			name varchar(100) NOT NULL,
       			password varchar(100) NOT NULL,
       			PRIMARY KEY (name)
       		);
       `

       INSERT_USER = `
               INSERT INTO TABLE users VALUES(?,?);
       `
       UPDATE_USER = `
               UPDATE TABLE users SET password=? WHERE name=?;
       `
       CREATE_MSG = `
        	CREATE TABLE IF NOT EXISTS messages (
       			sender varchar(100) NOT NULL,
       			recipient varchar(100) NOT NULL,
                        message varchar(1000),
                        contentmetadata varchar(50) NOT NULL,
                        timestamp timestamp NOT NULL,
       			PRIMARY KEY (sender,receipient)
       		);
       `

       INSERT_MSG = `
               INSERT INTO TABLE messages VALUES;
       `
)
