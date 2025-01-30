FROM node:22-alpine3.21 as nodeBuilder

WORKDIR /app

COPY package.json .
COPY package-lock.json .

RUN npm install

COPY styles.css .
COPY internal internal

RUN npm run build

FROM golang:1.23-alpine3.21 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/main cmd/api/main.go

FROM alpine:3.21

EXPOSE 80
WORKDIR /app

COPY --from=builder /app/main /app/main
COPY assets assets
COPY --from=nodeBuilder /app/assets/styles.css assets/styles.css

CMD ["/app/main", "-port", "80"]
