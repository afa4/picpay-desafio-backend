FROM golang:1.22 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./src

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /build/app ./app

ENTRYPOINT ["./app"]