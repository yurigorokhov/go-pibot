FROM golang:1.6-onbuild

EXPOSE 8080

# build application
RUN go get -d -v
RUN go build -v

