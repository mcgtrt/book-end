# book-end - Hotel Reservation API

## Project outline
- users -> book a hotel room
- admins -> check the hotel reservations
- Authentication and authorization -> JWT tokens
- Hotels -> CRUD API -> JSON
- Rooms -> CRUD API -> JSON
- Scripts -> database management -> seeding, migration


## Project environment variables
```
MONGO_DB_URI=mongodb://localhost:27017
MONGO_DB_NAME=book-end
TEST_MONGO_DB_NAME=book-end-test
JWT_SECRET=superstrongsecretkeyNOBODYknows
HTTPS_LISTEN_ADDRESS=:3000
```

## Resources
### Mongodb driver 
Documentation
```
https://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongodb client
```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber 
Documentation
```
https://gofiber.io
```

Installing gofiber
```
go get github.com/gofiber/fiber/v2
```

## Docker
### Installing mongodb as a Docker container
```
docker run --name mongodb -d mongo:latest -p 27017:27017
```

## GoDotEnv Load
```
go get github.com/joho/godotenv
```