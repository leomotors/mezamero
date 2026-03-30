# Build
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /mezamero ./cmd/mezamero

# Run — mount host config at /config/config.yaml
FROM scratch
COPY --from=build /mezamero /mezamero
EXPOSE 8080
ENTRYPOINT ["/mezamero"]
CMD ["-config", "/config/config.yaml", "-addr", ":8080"]
