FROM golang:1.19

WORKDIR /app

COPY ./go.mod ./go.sum /app
RUN go mod download && go mod verify

COPY . /app
RUN go build /auth.go

EXPOSE 5000

CMD ["go", "run", "/auth.go"]
