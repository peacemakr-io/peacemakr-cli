FROM peacemakr-cli-dependencies as builder


ADD go.mod go.sum /go/src/github.com/peacemakr/peacemakr-cli/

ADD . /go/src/github.com/peacemakr/peacemakr-cli/

ENV GOPATH=/go
WORKDIR /go/src/github.com/peacemakr/

RUN go install /go/src/github.com/peacemakr/peacemakr-cli/

WORKDIR /go/bin
ADD peacemakr.yml ./
FROM alpine

RUN apk update && apk upgrade && apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /go/src/github.com/peacemakr/peacemakr-cli/vendor/ /go/src/
COPY --from=builder /go/pkg/mod/github.com/peacemakr-io/peacemakr-go-sdk@v0.0.11-0.20201217045855-77efbe3bd32f/pkg/crypto/lib/libpeacemakr-core-crypto.so /lib/

WORKDIR /go/bin/
COPY --from=builder /go/bin/peacemakr-cli /go/bin/peacemakr-cli

CMD ./peacemakr-cli
