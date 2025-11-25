# Build stage: compile static Go binary on tiny Alpine builder
FROM golang:alpine AS build

# add git only for module fetches
RUN apk add --no-cache git

WORKDIR /app/proxy
COPY . /app

# download modules (harmless if not using modules)
RUN go mod download || true

ENV CGO_ENABLED=0
# static, trimmed binary suitable for scratch
RUN GOOS=linux go build -trimpath -ldflags="-s -w" -o /proxy .

# Optional: compress the binary with UPX (uncomment if you want smaller image,
# note: UPX can break some binaries or interfere with ptrace/debugging)
# RUN apk add --no-cache upx && upx --best /proxy

# Alpine stage: only install runtime packages so we can copy the small runtime files
FROM alpine:latest AS alpine-base
RUN apk add --no-cache tor tzdata ca-certificates

# Final minimal image
FROM scratch

# CA certs (single bundle file)
COPY --from=alpine-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy a single timezone file (UTC) instead of the whole zoneinfo tree.
# If your app needs a different timezone, replace 'UTC' with the specific zone you need.
COPY --from=alpine-base /usr/share/zoneinfo/UTC /etc/localtime

# Copy only Tor's geoip files (as requested)
COPY --from=alpine-base /usr/share/tor/geoip* /usr/share/tor/

# Application binary
COPY --from=build /proxy /bin/proxy
ENTRYPOINT [ "/bin/proxy" ]
