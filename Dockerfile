FROM golang:latest as builder
ADD . /go/src/github.com/Haelium/BloomReachTest
RUN go get ./...
RUN go install github.com/Haelium/BloomReachTest/cmd/userapi

ENTRYPOINT /go/bin/userapi
FROM golang:latest
EXPOSE 8080
COPY --from=builder /go/bin/userapi .
ADD run.sh .
RUN chmod +x run.sh
CMD ["./run.sh"]