FROM debian:buster

EXPOSE 3000

ENV GODEBUG netdns=go

ADD rotocopter /bin/
ENTRYPOINT ["/bin/rotocopter"]
