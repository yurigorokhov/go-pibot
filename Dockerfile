FROM golang:1.12

ADD . /go

EXPOSE 8080

# build application
RUN go get -d -v
RUN go build -v -o app
