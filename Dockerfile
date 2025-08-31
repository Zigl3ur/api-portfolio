FROM golang:1.24.6-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM golang:1.24.6-alpine AS runner

WORKDIR /app

COPY --from=base /app/main .

CMD ["/app/main"]