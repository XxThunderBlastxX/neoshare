FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY neoshare /usr/bin/neoshare

EXPOSE 8080

ENTRYPOINT ["/usr/bin/neoshare"]