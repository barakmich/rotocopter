FROM golang:1.15

RUN go build -o rotocopter .

EXPOSE 3000

ENTRYPOINT ["/bin/rotocopter"]
