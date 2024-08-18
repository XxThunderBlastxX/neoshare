FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY ./static /usr/bin/static
COPY neoshare /usr/bin/neoshare

EXPOSE 8080

ENTRYPOINT ["/usr/bin/neoshare"]