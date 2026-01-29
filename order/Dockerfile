FROM golang:1.25 AS builder

WORKDIR /app

ENV GOPROXY=direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/main.go


FROM scratch

COPY --from=builder /app/server /server


EXPOSE 3000

ENTRYPOINT ["/server"]