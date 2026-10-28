FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /go-api ./cmd

FROM alpine:latest
WORKDIR /
COPY --from=builder /go-api /go-api
COPY --from=builder /app/internal/data/residential_buildings_bcn.geojson /internal/data/residential_buildings_bcn.geojson
COPY --from=builder /app/internal/data/sensor_building_map_cassandra.json /internal/data/sensor_building_map_cassandra.json
#COPY --from=builder /app/internal/data/sensor_building_map.json internal/data/sensor_building_map.json
EXPOSE 8090
CMD ["/go-api"]