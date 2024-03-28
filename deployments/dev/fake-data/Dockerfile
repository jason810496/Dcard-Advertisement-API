# build stage
FROM golang:alpine AS builder
ADD . /src
RUN cd /src && go build -v -o bin/fake-data ./cmd/fake-data

ADD .env /src/.env

# final stage
FROM alpine
WORKDIR /app
COPY --from=builder /src/bin/fake-data /app/
COPY --from=builder /src/.env /app/.env

EXPOSE 8080
CMD ["./fake-data"] 