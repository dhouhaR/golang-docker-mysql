FROM golang:alpine

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /golang-docker-todo

ADD . .

RUN go mod download

#todo : test that
#RUN go get github.com/githubnemo/CompileDaemon
#ENTRYPOINT CompileDaemon  -command="./golang-docker-todo"
ENTRYPOINT go build  && ./golang-docker-todo

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./main"]