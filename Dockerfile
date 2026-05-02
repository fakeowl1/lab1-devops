FROM golang:1.26-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/app ./cmd/mywebapp/main.go

FROM gcr.io/distroless/static
COPY --from=builder /usr/local/bin/app /
COPY --from=builder /usr/src/app/templates /templates

ENTRYPOINT ["/app"]
