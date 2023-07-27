FROM golang:latest

ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get install -y postgresql-client

RUN chmod +x wait-for-postgres.sh

COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN ln -s /go/bin/linux_amd64/migrate /usr/local/bin/migrate

RUN go build -o ./bug-tracker ./cmd/bug-tracker/main.go

EXPOSE 7000

CMD [ "./bug-tracker" ]