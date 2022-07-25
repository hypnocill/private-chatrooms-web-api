FROM golang:1.14-alpine

WORKDIR /api

COPY . /api

## Add this go mod download command to pull in any dependencies
RUN go mod download

RUN go build -o main .
## Our start command which kicks off
## our newly created binary executable
CMD ["/api/main"]

EXPOSE 5000
