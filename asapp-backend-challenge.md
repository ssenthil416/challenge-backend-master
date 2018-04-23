ASAPP Backend Challenge
=======================

Welcome to your challenge project!

For this challenge, we ask that you implement a solution at home in your own time. When you're ready, please send us your results - including code, schema definitions, commit history if available, and (if you're not using the starter kit) a readme with setup instructions.


The Details
-----------

Your challenge is to design and implement a basic chat backend. While building a frontend isn't part of this challenge, it helps to keep in mind how a client would interact with your server and what the experience will be like for the end user.

Your API should be exposed in a format that can be called from a web app or mobile client - for example, JSON and HTTP. Your server should support the following requests:

* Create User
	
Takes a username and password and creates a new user in a persisted data store. For example, the endpoint might accept PUT or POST at /users.

* Send Message

Takes a sender, recipient, and message and saves that to the data store. For example, the endpoint might accept PUT or POST at /messages. Three different message types are supported. (1) is a basic text-only message. (2) is an image link. (3) is a video link. The image and video links are saved with some additional metadata: width and height for the image, length of the video and source (YouTube, Vevo) for the video. For the purpose of the challenge, you can hard-code the metadata to some fixed values when you're saving the message.

* Fetch Messages

Takes two users and loads all messages sent between them. This call should also take two optional parameters in order to support pagination: the number of message to show per page and which page to load. For example, the endpoint might accept GET at /messages.



Starter Kit
-----------
We've included a Docker-based starter kit with some commonly used languages. The kit includes a skeleton db and server code. You're welcome but not required to use this as the basis of your project. To use the starter kit, follow the readme in challenge-eng-base-master.

If opting not to use the starter kit, Please include a script that builds and launches your server.


Suggestions
-----------

* We very much value code quality and technical design. Think about the structure of your APIs, your data models, and the readability of your code.
* At ASAPP, we use a lot of Go and MySQL. For the challenge, we’d like you to be able to work in languages with which you’re comfortable, but we do suggest the following:
```
    Backend: Go, Python, Java, Node
    Database: SQL (including SQLite, MySQL, Postgres)
```
* Use open source libraries rather than reinventing the wheel. Here are a couple of relevant tools that we use:
```
    github.com/go-sql-driver/mysql
    golang.org/x/crypto/bcrypt
```
* Please include a sample request (cURL commands, Postman collection, etc) for each of your API endpoints.
* Please don't use the trademark ASAPP in the project. We hope the project is work that you're proud of, and we want you to be able to share it with others or make it public should you wish to.
* Have fun!


Follow-up discussion
--------------------
We’ll discuss this as part of the project review. Don't worry if you don't have all the answers off the top of your head. We’re very much looking for your ability to reason about and work through these kinds of questions.

How well does your project scale? What if the number of users grow to 1000? To 1000000? (And the conversation history grows too.)

What if you had users around the globe? How do you keep the application responsive? (Latency becomes problematic if you’re still running in just one region. But if you have servers in Japan and servers in the US, how do they coordinate?)
