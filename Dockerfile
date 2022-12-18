FROM golang:1.17-bullseye as server
WORKDIR /src

COPY . .
RUN go get -v -t -d . && go build -o bin/checker ./main.go
ENTRYPOINT go run main.go