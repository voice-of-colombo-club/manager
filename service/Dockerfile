FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /voc-manager

EXPOSE 8080

CMD [ "/voc-manager" ]