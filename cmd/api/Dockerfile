FROM golang:1.23-alpine AS build-stage

WORKDIR /app

# NOTE: download dependencies before copying the source code to take advantage of Docker's layer caching
# ref: https://docs.docker.com/build/guide/layers/#cached-layers
COPY ../../go.mod ./
RUN go mod tidy

COPY ../.. .

# Build for ARM64 architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o expense-tracker-go ./cmd/api/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/expense-tracker-go ./expense-tracker-go
COPY --from=build-stage /app/migrations ./migrations/

USER nonroot:nonroot

CMD ["./expense-tracker-go"]