FROM golang:1.17.1-alpine3.14 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY limits.conf /etc/security/limits.conf

WORKDIR /build
COPY go.mod go.sum main.go ./

RUN go mod download

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

FROM scratch

COPY --from=builder /dist/main .

ENTRYPOINT ["/main"]