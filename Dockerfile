FROM peacemakr-cli-dependencies as builder


ADD go.mod go.sum /go/src/github.com/peacemakr/peacemakr-cli/

ADD main.go /go/src/github.com/peacemakr/peacemakr-cli/main.go

ENV GOPATH=/go
WORKDIR /go/src/github.com/peacemakr/

RUN go install /go/src/github.com/peacemakr/peacemakr-cli/

WORKDIR /go/bin
ADD peacemakr.yml ./
# FROM alpine

# RUN apk update && apk upgrade && apk add --no-cache ca-certificates && update-ca-certificates

# COPY --from=builder /go/src/github.com/peacemakr/peacemakr-cli/vendor/ /go/src/

# WORKDIR /go/bin/
# COPY --from=builder /go/bin/peacemakr-cli /go/bin/peacemakr-cli

CMD ./peacemakr-cli
