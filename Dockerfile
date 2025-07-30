FROM golang:1.24.4-alpine AS builder

WORKDIR /go/src/masstrack

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s" -o masstrack /go/src/masstrack/main.go

FROM alpine:3.22

ENV GIN_MODE=release

COPY --from=builder /go/src/masstrack/masstrack /usr/bin/masstrack

EXPOSE 8080

ENTRYPOINT ["/usr/bin/masstrack"]
