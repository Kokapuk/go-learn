FROM golang:1.26.5-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /api \
    ./cmd/api


FROM scratch

COPY --from=build /api /api
COPY --from=build /etc/ssl/certs/ca-certificates.crt \
    /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/api"]