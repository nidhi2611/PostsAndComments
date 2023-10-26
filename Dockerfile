FROM golang:latest

EXPOSE 3000

COPY . /app

WORKDIR /app
RUN go build
ENTRYPOINT ["/app/goassignment-postandcomments"]