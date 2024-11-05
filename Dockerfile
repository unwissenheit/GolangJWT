FROM golang:1.23-alpine

WORKDIR /app

CMD [ "sh", "entrypoint.sh" ]