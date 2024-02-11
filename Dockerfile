FROM golang:1.22.0 as build

WORKDIR /app

# Copy golang dependency manifests
COPY go.mod .
COPY go.sum .

# Cache the downloaded dependency in the layer
RUN go mod download

# Copy the source code
COPY . .

RUN go build -o /app/dns_server .

# COPY --from=build /app/dns_server /bin/dns_server
# ENTRYPOINT ["/bin/dns_server"]
ENTRYPOINT ["/app/dns_server"]
