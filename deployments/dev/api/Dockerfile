# build stage
FROM golang:1.21-alpine AS builder
ADD . /src
RUN cd /src && go mod download
RUN cd /src && go build -v -o bin/api ./cmd/api

ADD .env /src/.env

# final stage
FROM alpine
WORKDIR /app
COPY --from=builder /src/bin/api /app/
COPY --from=builder /src/.env /app/.env

EXPOSE 8080
CMD ["./api"] 