FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY ./static ./static
COPY neoshare ./neoshare

EXPOSE 8080

ENTRYPOINT ["./neoshare"]