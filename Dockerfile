FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY ./static .
COPY neoshare .

EXPOSE 8080

ENTRYPOINT ["neoshare"]