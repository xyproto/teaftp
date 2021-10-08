FROM golang:1.17 as builder
LABEL maintainer="xyproto@archlinux.org"

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags "-s" -a -v

FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/teaftp /usr/bin/teaftp

COPY --from=builder /app/static /srv/tftp
WORKDIR /srv/tftp

ENV PORT 69
EXPOSE 69

CMD /usr/bin/teaftp
