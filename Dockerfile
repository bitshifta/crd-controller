FROM golang:1.10.3-alpine3.7 as builder

COPY . /go/src/github.com/bitshifta/crd-controller/

RUN cd /go/src/github.com/bitshifta/crd-controller/ && \
    go build *.go

FROM alpine:latest

COPY --from=builder /go/src/github.com/bitshifta/crd-controller/controller .

CMD ["./controller"]

