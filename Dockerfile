FROM golang:1.15 as go

RUN mkdir /rotocopter
ADD . /rotocopter/
WORKDIR /rotocopter

RUN go build -o rotocopter .

FROM debian:buster

COPY --from=go /rotocopter/rotocopter /bin
EXPOSE 3000

ENTRYPOINT ["/bin/rotocopter"]
