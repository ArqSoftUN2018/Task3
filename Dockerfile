FROM golang:1.10.4
RUN mkdir /app
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]
<<<<<<< HEAD
=======

>>>>>>> test
EXPOSE 3002
