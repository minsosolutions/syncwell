FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o /mockapi ./cmd/mockapi

FROM alpine:latest
WORKDIR /app
COPY --from=build /mockapi /usr/local/bin/mockapi
EXPOSE 8080
CMD ["mockapi", "-data", "/app/data", "-config", "/app/config/syncwell.json", "-addr", ":8080"]
