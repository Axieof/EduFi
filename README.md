# EduFi

EduFi is a University Learning Management System devloped as 23 seperate microservices for a class assignment. 
Each student in the class will do a microservice to that we all will containerize with docker and deploy onto a server, 
which then allows the services to talk to each other. The purpose of the assignment is to be able to work with others to 
develop microservices that can work independently and communicate with one another.

# Setup Instructions
Each Service requires the Go [Echo Framework](https://echo.labstack.com/guide/) package, 
with the database service needing the mongo driver package

# Microservice Network Diagram
![EduFi Architecture](https://github.com/Axieof/EduFi/blob/master/setup/EduFI_Network_Diagram_Final_Updated.png)

For my contribution to EduFi, I have worked on the Marks Entry function of the Learning Management System.
The marks Entry portion has the following requirements
- List all Modules tought by Tutor
- List all Classes tought by Tutor
- List all Students tought by Tutor
- View Marks of Students
- Update Marks of Students
- Deposit Tokens of Equivalent Value into Student's Wallet

As EduFi is made of up many microservices working together, my microservice does act as a start point. 
Instead, a user would log in as a tutor through the Authentication service and select marks entry from
a navigation bar in order to access my service. 

1. Upon accessing my service, the authentication service 
passes the tutor's Identification number to my Frontend service.

2. My Frontend then sends the Tutor ID to the Tutor Service to receive information on the tutor such as
modules, classes and students taught. 

3. The Frontend service then uses the information received on the tutor to present the user with a
User Interface displaying the modules, classes and students taught. 

4. The Tutor can then enter marks for a student from a class they teach and submit the mark.

5. The mark then gets submtited back to the frontend service, which then adds the mark submitted along with
the tutor's and student's information to the Marks service.

6. The Marks service then processes the information and sends the information to 2 services
- The Marks service sends the information to the Database service to be added into the Marks Entry service
- The Marks service sends the information to the Wallet service to provide the student with tokens corresponding to module

7. The Database service then calculates the current semester date and sends the information to the Mongo database to be stored within a colelction named after the current semester.

8. All services send a erply according to the service that calls them, to confirm that the mark has been added to the database, as well as sent to the Wallet microservice.

# Docker Links

Marks Service: [Docker Marks Image](https://hub.docker.com/repository/docker/axieof/edufi-marks)
	
Database Service: [Docker Database Image](https://hub.docker.com/repository/docker/axieof/edufi-database)
	
Frontend Service: [Docker Frontend Image](https://hub.docker.com/repository/docker/axieof/edufi-frontend)

# Setup

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
This route is to be used with a curl command to 
receive a reply if the server is up and running.
```

Example Request
CURL command
```
curl http://localhost:8121/api/V1/checkapi
```

- POST (http://localhost:8121/api/V1/marksSubmit/:studentID)
```
This route is used when the Frontend service posts marks of a 
student entered by a tutor, to be processed
by the Marks Service.
```

### Frontend Service
- GET (http://localhost:8120/checkapi)
```
This route is to be used with a curl command to receive a 
reply if the server is up and running.
```

Example Request
CURL command
```
curl http://localhost:8120/checkapi
```

- GET (http://localhost:8120/marksDashboard/:tutorID)
```
This route is when the user attempts accesses the marks 
dashboard, with their ID passed as a query parameter.
The user will then be greeted with a UI screen of the Modules, 
Classes and Students that they teach, and then they can provide
marks to each student.
```
- POST (http://localhost:8120/marksEntry/:studentID)
```
This route is when the tutor submits a mark of a student, which 
then posts to this route within the Frontend service to be processed 
and sent to the Marks service.
```

### Database Service
- GET (http://localhost:8129/api/V1/checkapi)
```
This route is to be used with a curl command to receive a 
reply if the server is up and running.
```

Example Request
CURL command
```
curl http://localhost:8129/api/V1/checkapi
```

- POST (http://localhost:8129/api/V1/database/postMarks)
```
This route is used when the database service has received marks 
of a student from the Student service, which then gets processed 
and sent of to the mongo database container
```
