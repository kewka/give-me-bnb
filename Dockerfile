FROM golang:1.17-alpine3.14 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED 0
RUN go build .

FROM alpine:3.14
WORKDIR /app
RUN apk add chromium tor
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh
COPY --from=builder /app/give-me-bnb .
ENV PATH="${PWD}:${PATH}"
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["give-me-bnb"]