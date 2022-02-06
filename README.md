# EduFi

EduFi is a University Learning Management System devloped as 23 seperate microservices for a class assignment. 
Each student in the class will do a microservice to that we all will containerize with docker and deploy onto a server, 
which then allows the services to talk to each other. The purpose of the assignment is to be able to work with others to 
develop microservices that can work independently and communicate with one another.

# Setup Instructions
Each Service requires the Go [Echo Framework](https://echo.labstack.com/guide/) package, 
with the database service needing the mongo driver package

# Microservice Network Diagram

FrontEnd Service Setup Commands

```
cd src/Frontend
go get "github.com/labstack/echo/v4"
go run main.go
```

Marks Service Setup Commands

```
cd src/Marks
go get "github.com/labstack/echo/v4"
go run main.go
```

Database Service Setup Commands

```
cd src/Database
go get "github.com/labstack/echo/v4"
go get "go.mongodb.org/mongo-driver/bson"
go get "go.mongodb.org/mongo-driver/bson/primitive"
go get "go.mongodb.org/mongo-driver/mongo"
go get "go.mongodb.org/mongo-driver/mongo/options"
go run main.go
```

# API Endpoints
### Marks Service
- GET (http://localhost:8121/api/V1/checkapi)
```
This route is to be used with a curl command to receive a reply if the server is up and running.
```
- POST (http://localhost:8121/api/V1/marksSubmit/:studentID)
```
This route is used when the Frontend service posts marks of a student entered by a tutor, to be processed
by the Marks Service.
```
### Frontend Service
- GET (http://localhost:8120/checkapi)
```
This route is to be used with a curl command to receive a reply if the server is up and running.
```
- GET (http://localhost:8120/marksDashboard/:tutorID)
```
This route is when the user attempts accesses the marks dashboard, with their ID passed as a query parameter.
The user will then be greeted with a UI screen of the Modules, Classes and Students that they teach, and then they can provide marks to each student.
```
- POST (http://localhost:8120/marksEntry/:studentID)
```
This route is when the tutor submits a mark of a student, which then posts to this route within the Frontend service to be processed and sent to the Marks service.
```
### Database Service
- GET (http://localhost:8129/api/V1/checkapi)
```
This route is to be used with a curl command to receive a reply if the server is up and running.
```
- POST (http://localhost:8129/api/V1/database/postMarks)
```
This route is used when the database service has received marks of a student from the Student service, which then gets processed and sent of to the mongo database container
```
