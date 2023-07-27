# Bug Tracker

Simple bug tracker written in golang.
Technologies:
    + PostgreSQL
    + Redis
    + Kafka

## How to run

1. Clone git repo
2. Go to [mail-sender](https://github.com/samuraivf/mail-sender) and read README.
3. Create `.env` file with these environment variables:
```
POSTGRES_PASSWORD=...
POSTGRES_URL=...
MAIL_FROM=...
MAIL_FROM_APP_PASSWORD=...
SMTP_HOST=...
SMTP_PORT=...
BUG_TRACKER_ADDRESS=...
```
instead of `...` there should be your data.
4. Build bug-tracker Docker image:
``` bash
$ docker build -t bug-tracker .
```
5. Run:
``` bash
$ docker-compose build && docker-compose up
```