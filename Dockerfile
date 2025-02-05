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

RUN go install github.com/a-h/templ/cmd/templ@v0.3.819

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN templ generate && go build -o /app/main cmd/api/main.go

FROM alpine:3.21

EXPOSE 80
VOLUME [ "/data" ]

WORKDIR /app

RUN mkdir -p /app/assets /data
RUN \
	wget -O /app/assets/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js \
	&& wget -O /app/assets/alpine.js https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js

COPY ./assets/favicon.svg assets/favicon.svg

COPY --from=nodeBuilder /app/assets/styles.css assets/styles.css

COPY --from=builder /app/main /app/main

CMD ["/app/main", "-port", "80", "-itemsJsonPath", "/data/items.json"]
