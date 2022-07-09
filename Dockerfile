FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY config/ ./config
COPY sql/ ./sql
COPY *.go ./

RUN go build -o /kudo-giver

ENV ENVIRONMENT=pre
EXPOSE 8080
CMD [ "/kudo-giver" ]
