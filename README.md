# Users Service

This is a small micro-service written in Go, using a containarized api and database which provides functionalities to add, delete, update and read users with the following fields: first name, last name, nickname, password, email and country.

### How to run?

### Add a new user

> POST /users

### Edit user

> PUT /users/:userid

### Remove user

> DELETE /users/:userid

### Get users

> GET /users/:userid

This route accepts the following parameters: ...
