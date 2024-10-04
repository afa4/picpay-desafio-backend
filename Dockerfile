FROM golang:1.22 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./src .

RUN go build -o /bin/app

ENTRYPOINT ["/bin/app"]