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
