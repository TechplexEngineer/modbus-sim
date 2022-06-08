# ref: https://github.com/jeremyhuiskamp/golang-docker-scratch
################################
# STEP 1 build executable binary
################################
FROM golang:alpine as builder

ARG GOARCH
ENV GOARCH ${GOARCH}
ENV CGO_ENABLED=0

WORKDIR /src

COPY . .

RUN go env

# Static build required so that we can safely copy the binary over.
# `-tags timetzdata` embeds zone info from the "time/tzdata" package.
RUN go build -ldflags '-extldflags "-static"' -tags timetzdata -o modbus ./...

################################
# STEP 2 build a small image
################################
FROM scratch

ENV TZ "America/New_York"

# ca-certificates to allow secure connections to other https servers
# NB: this pulls directly from the upstream image, which already has ca-certificates:
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

# Copy our static executable.
COPY --from=builder /src/modbus /app/modbus

EXPOSE 1502

ENTRYPOINT ["/app/modbus"]