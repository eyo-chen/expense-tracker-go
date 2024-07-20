FROM golang:1.22.2-alpine AS build-stage

WORKDIR /app

# NOTE: download dependencies before copying the source code to take advantage of Docker's layer caching
# ref: https://docs.docker.com/build/guide/layers/#cached-layers
COPY go.mod ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o expense-tracker-go ./cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/expense-tracker-go ./expense-tracker-go
COPY --from=build-stage /app/.env ./.env
COPY --from=build-stage /app/migrations ./migrations/

USER nonroot:nonroot

CMD ["./expense-tracker-go"]