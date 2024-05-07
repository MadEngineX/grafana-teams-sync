FROM golang:1.21.3 AS dependency

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy

FROM dependency as build

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o /out/grafana-teams-sync cmd/grafana-teams-sync/main.go

FROM alpine:3.19.1 AS runtime

WORKDIR /app

COPY --from=build /out/grafana-teams-sync /app/

ENTRYPOINT ["/app/grafana-teams-sync"]
