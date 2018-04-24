package main

const (
	CREATE_USERS = `
        	CREATE TABLE IF NOT EXISTS users (
       			name varchar(100) NOT NULL,
       			password varchar(100) NOT NULL,
       			PRIMARY KEY (name)
       		);
       `

	INSERT_USER = `
               INSERT users SET name=?, password=?;
       `
	UPDATE_USER = `
               UPDATE  users SET password=? WHERE name=?;
       `
	CREATE_MSG = `
        	CREATE TABLE IF NOT EXISTS messages (
       			sender varchar(100) NOT NULL,
       			recipient varchar(100) NOT NULL,
                        message varchar(1000),
                        contentmetadata varchar(50) NOT NULL,
                        timestamp timestamp NOT NULL,
       			PRIMARY KEY (sender,recipient)
       		);
       `

	INSERT_MSG = `
               INSERT INTO TABLE messages VALUES;
       `
)
