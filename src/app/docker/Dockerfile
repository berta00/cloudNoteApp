FROM golang:1.19

WORKDIR /app

COPY . /app
RUN go mod download && go mod verify
RUN go build /auth.go

EXPOSE 5000

CMD ["go", "run", "/auth.go"]
