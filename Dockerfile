FROM golang:1.25.7-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/server/main.go

FROM golang:1.25.7-alpine AS runner

WORKDIR /app

COPY --from=base /app/main .

ENV APP_ENV=production

CMD ["/app/main"]