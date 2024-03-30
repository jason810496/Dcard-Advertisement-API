# build stage
FROM golang:1.21-alpine AS builder
ADD . /src
RUN cd /src && go mod download
RUN cd /src && go build -v -o bin/scheduler ./cmd/scheduler

ADD .env /src/.env

# final stage
FROM alpine
WORKDIR /app
COPY --from=builder /src/bin/scheduler /app/
COPY --from=builder /src/.env /app/.env

CMD [ "./scheduler" ] 