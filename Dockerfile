FROM alpine:3.6 as alpine

RUN apk add -U --no-cache ca-certificates

EXPOSE 3000

ENV GODEBUG netdns=go

ADD rotocopter /bin/
ENTRYPOINT ["/bin/rotocopter"]
