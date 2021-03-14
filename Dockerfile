FROM golang:1.16-alpine
WORKDIR /var/app
COPY . .
RUN go build -o bin .

FROM alpine
WORKDIR /var/app
COPY --from=0 /var/app/bin /var/app/bin
RUN apk add ca-certificates
CMD ./bin
