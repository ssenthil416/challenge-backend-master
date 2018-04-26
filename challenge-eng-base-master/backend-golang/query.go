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
                        timestamp varchar(100) NOT NULL,
       			PRIMARY KEY (sender,recipient,timestamp)
       		);
        `

	INSERT_MSG = `
               INSERT messages SET sender=?, recipient=?, message=?, contentmetadata=?, timestamp=?;
        `

	UPDATE_MSG = `
               UPDATE  messages SET message=? WHERE sender=? and recipient=? and timestamp=?;
        `

	DROP_MSG = `DROP TABLE messages;`

	SELECT_MSG = `
               SELECT sender, recipient, message, contentmetadata, timestamp FROM messages WHERE sender=? AND recipient=?;
        `
)
