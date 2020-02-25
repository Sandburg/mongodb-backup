FROM golang:latest AS builder

WORKDIR /usr/src

COPY . .

RUN cd src/ && go get
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /go/bin/backup ./src/*.go

FROM mongo:latest

COPY --from=builder /go/bin/backup /backup

WORKDIR /