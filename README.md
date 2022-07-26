# Private Chatrooms Web API

## How to start the project on your machine

 1. Clone the project
```
$ git clone https://github.com/hypnocill/private-chatrooms-web-api.git
```

2. Start docker compose inside the folder of the project.
```
$ docker-compose up
```

3. Opening http://localhost:5000/ should show "Private Chatrooms v1" text if the project is running

4. Create room (make sure to use the received params from the 'response' parts below in the example requests)
```
$ curl -X POST http://localhost:5000/room/create  -H 'Content-Type: application/json'  -d '{"password":"TEST_PASSWORD"}';
```
response:
```
{"Id":"1658790109809"}
```

5. Authenticate for the created room (Take the ID from the response of the previous step)
```
$ curl -X POST http://localhost:5000/room/auth  -H 'Content-Type: application/json'  -d '{"roomId":"1658790109809","password":"TEST_PASSWORD"}';
```
response:
```
{"Token":"5fd924625f6ab16a19cc9807c7c506ae1813490e4ba675f843d5a10e0baacdb8"}
```

6. Join the room with a websocket client of your preference with url ws://localhost:5000/room/join/1658790109809/YOUR_USER_NAME
(Google Chrome extension - Simple Websocket Client is my preference)

The first message sent to the room should be the authentication token with the following string
```
{"PayloadType":0,"Message":"5fd924625f6ab16a19cc9807c7c506ae1813490e4ba675f843d5a10e0baacdb8"}
```
After that, connection and authentication for the room should be established and normal messages can be sent.
Each user will be notified when a new user joins the room
See '/service/chatroom/payload_types.go' for the different types of payloads and how they are used in the project.

![Alt text](/res/chat.png?raw=true "Example chat with Simple Websocket Client")


