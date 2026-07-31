FROM golang:1.26 AS build
WORKDIR /app/src
COPY src/go.mod ./
RUN go mod download
COPY src/ ./
RUN go build -o /tinycolor-compat ./cmd/tinycolor-compat

FROM scratch
COPY --from=build /tinycolor-compat /tinycolor-compat
ENTRYPOINT ["/tinycolor-compat"]
