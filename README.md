# Users Service
 
This is a small micro-service written in Go, using a containarized api and database which provide functionalities to add, delete, update and read users with the following fields: first name, last name, nickname, password, email and country.
Also there is a status route to check service health.

### How to run?
API + MONGO DB
> docker-compose build && docker-compose up


Tests (from route dir)
> gotest ./...

The API runs on localhost on port 8080.

### Assumptions made during development
- Regarding Mongo Adapters, I decided to separate Mongo Adapter for two simple reasons: the first is makes easier to test, the `mongo.go` doesn't have external libraries so coding adapters this way makes mocking pretty simple, also there is no point on testing functions from external libraries because libs should have their own tests. The second reason is it makes easier if we decide to move our database provider, (e.g) if we need to change for instance `mongo.db` to `mysql`, mongo is wrapped so we wouldn't need to change our API code.

- I didn't go too deep on the unit testing. The idea in this microservice is to provide the logic on how I write tests and not test all the possible cases. I am a table driven tests fan cause because imho they are easier to maintain.

- main.go isn't tested and I assume if things go wrong in the main.go the service shouldn't start.

- The GET user routes checks if the query parameter is expected to avoid SQL Injections, also for the purpose of this service, assume this route only allows one single parameter per key, i.e, you can't provide twice the same parameter.

- I have used `guid` instead of Mongo ObjectIDs to represent user ids - I am just used to it.  

### Possible extensions or improvements to the service

- Move the health route to a separate file and possibly check more stuff alongside the database.
- Unit tests should cover all the possible scenarios.
- I'd write e2e tests, by running the api on a test container and making calls to the api then making assertions to responses and database documents etc.
- I would add pagination to the GET users route, because if you have millions of users you can't just get them all at a single time.
- I'd Hash & Salt the passwords.
- Body fields validation, e.g, check email format or check if country is valid.

## API Endpoints 

### Add a new user

> POST /users

body: 
```
    {
    "nickname": "jpaldi",
    "email": "jpaldi@email.pt",
    "first_name": "joao",
    "last_name": "aldi",
    "password": "S3CR3T",
    "country": "PT"
}
```

If the User is successfully created the service returns a 200 Status Code and returns the new user including the new id.
If fields are missing the service returns a 400 Status Code and reports the errors.

### Edit user

> PUT /users/:userid

body: 
```
    {
    "nickname": "jpaldi",
    "email": "jpaldi@email.pt",
    "first_name": "joao",
    "last_name": "aldi",
    "password": "S3CR3T",
    "country": "PT"
}
```
If the User is successfully created the service returns a 200 Status Code and returns the updated document for this user.
If fields are missing the service returns a 400 Status Code and reports the errors.

### Remove user

> DELETE /users/:userid

If the User is successfully deleted the service returns a 200 Status Code

### Get users

> GET /users?country=PT

This route accepts the following parameters: nickname, first_name, country, last_name, email.

Response:
Status Code 200
body: 
```
[
        {
            "nickname": "jpaldi",
            "email": "jpaldi@email.pt",
            "first_name": "joao",
            "last_name": "aldi",
            "password": "S3CR3T",
            "country": "PT"
        },
         {
            "nickname": "jpaldi2",
            "email": "jpaldi2@email.pt",
            "first_name": "joao2",
            "last_name": "aldi2",
            "password": "S3CR3T",
            "country": "PT"
        },
        ...
    ]
}
```
