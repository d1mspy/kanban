FROM golang:1.23-bookworm AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /app/server /app/cmd/main.go

FROM ubuntu:22.04
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/db/migrations ./db/migrations
CMD [ "./server" ]
