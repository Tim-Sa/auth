FROM golang:1.21.6-alpine AS builder

COPY . /github.com/Tim-Sa/auth/source/
WORKDIR /github.com/Tim-Sa/auth/source/

RUN go mod download
RUN go build -o ./bin/auth cmd/main.go

FROM alpine:3.19.1

WORKDIR /root/
COPY --from=builder /github.com/Tim-Sa/auth/source/bin/auth .

CMD [ "./auth" ]