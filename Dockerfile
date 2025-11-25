# Build stage: compile the Go binary on Alpine
FROM golang:alpine AS build

RUN apk add --no-cache git

WORKDIR /app/proxy
COPY . /app

# Prefer module-aware download if present; harmless otherwise
RUN go mod download || true

ENV CGO_ENABLED=0
# static linux binary suitable for scratch
RUN GOOS=linux go build -a -installsuffix cgo -ldflags='-s -w' -o /bin/proxy .

# Alpine stage: provide tor runtime files, tzdata and CA certs
FROM alpine:latest AS alpine-base
RUN apk add --no-cache tor tzdata ca-certificates

# Final minimal image
FROM scratch

# Certificates
COPY --from=alpine-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Timezone data
COPY --from=alpine-base /usr/share/zoneinfo /usr/share/zoneinfo

# Copy only Tor's geoip files (as requested)
COPY --from=alpine-base /usr/share/tor/geoip* /usr/share/tor/

# Application binary
COPY --from=build /bin/proxy /bin/proxy

ENTRYPOINT [ "/bin/proxy" ]

LABEL io.containers.autoupdate=registry
