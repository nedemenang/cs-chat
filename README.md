# CUSTOMER SERVICE CHAT

## PROJECT DESCRIPTION

This is a simple application that simulates a customer service chat. It allows for two types of roles; users and admin.

### Users can

- Register
- Login
- View their messages
- Send new messages and recieve feedback from the other party

### Admins can

- Register
- Login
- View their messages
- Send messages and recieve feedback from the other party
- View messages posted by users that have not been responded to by other admins

This project uses a number of technologies, some of which include

- ReactJs
- Golang
- Docker and Docker-compose
- Websockets
- MongoDB

## RUNNING THE PROJECT

To run the project you should have `Docker` and `Docker-Compose` installed on your system. In the `docker-compose.yml` insert your mongodb connection string in the `MONGODB_URI` environment variable as well as a random string for your `ACCESS_SECRET`. Open a terminal and navigate to the root directory of the project. Run the following command:

- `docker-compose up --build`

Open a browser and place this address on your address bar:

- `http://localhost:3000/`

## RUNNING TEST

### Backend Test

To run the backend test, navigate into the `chat-api` folder and run `go test ./... --cover`
