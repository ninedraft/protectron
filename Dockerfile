FROM golang:1.10 AS builder 
ENV TOKEN="" 
COPY . $GOPATH/src/github.com/ninedraft/protectron
RUN go install github.com/ninedraft/protectron/cmd/protectron
RUN echo $PATH
FROM alpine:3.6
COPY --from=builder /go/bin /go/bin
RUN export PATH=$PATH:/go/bin && \
    echo $PATH && \
    ls -l /go/bin
ENTRYPOINT ["/go/bin/protectron"]