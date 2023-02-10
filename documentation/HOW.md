# golang-docker-todo

This is the API part of the TO DO-IT application.
It contains a set of webservices that handle the two entities : **User** and **Task**.

For more details, there is a postman collection : _TODO_IT.postman_collection.json_

## Prerequisites Environement

To use this app, you need to have
    [Docker](https://docs.docker.com/get-docker/)
    &&
    [GoLang](https://go.dev/doc/install)
installed on your system.

## Build & Run APP

Clone source code from : [golang-docker-todo](https://github.com/dhouhaR/golang-docker-todo.git)

From command line, under the _**golang-docker-todo**_ folder execute

* docker-compose build
* docker-compose up

* OR docker-compose up --build

You are ready now to use the API.

From the navigator visit  : <http://localhost:8080>

[localhost](localhost_8080.png)

phpMyadmin is accessible from : <http://localhost:5003/>

## Examples

There is a pre-embedded database.
You can start by retreiving existant users, to use one of them for testing other APIs.

I'll take **dhouha@maisonduweb.com/dhouha123** account, to explain how it works.

1. First of all, you need to generate a token using, POST /login route.
by using json input for body param :
``{
    "email": "dhouha@maisonduweb.com",
    "password": "dhouha123"
}``

Sample result : "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NzYwMzc0MTQsInVzZXJfaWQiOjN9.wLvCIs8C_OHdLzZKWL6SVKSiOhjt7xrBwm6arWxcC3Q"

2. Use this token to retreive all tasks created by this account:
GET /Tasks

3. For more examples you could import the postman collection : [TODO-IT-Collection](TODO-IT.postman_collection.json)

## Shutting down APP

After successful testing of the application, you can shut it down using:

* docker-compose down

Or to remove volumes also,

* docker-compose down --remove-orphans --volumes

* docker system prune : to remove dangling image

## TODO LIST

Unit tests are not finished, I'll update the repository later.
Status and categories are not yet implemented.
Also some filters on retreiving Tasks are not yet implemented.
Improve json output (ex : hide crypted password)
