FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /orders-api ./cmd/server
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /orders-api /orders-api
LABEL org.opencontainers.image.source="https://github.com/bhuang02/orders-api"
EXPOSE 8080
ENTRYPOINT ["/orders-api"]
