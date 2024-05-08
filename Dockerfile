FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app


RUN go build -o theater-app ./cmd/main.go
CMD ["./theater-app"]