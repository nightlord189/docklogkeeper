FROM golang:1.21-alpine3.18 AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o ./cmd/app/main ./cmd/app

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.18.0

RUN apk --no-cache add tzdata

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /build/cmd/app/main /
COPY --from=builder /build/configs /configs
COPY --from=builder /build/static /static

EXPOSE 3010

ENTRYPOINT ["/main"]