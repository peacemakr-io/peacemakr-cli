FROM peacemakr-dependencies as builder

ADD main.go /go/src/peacemakr/peacemakr-cli/main.go

ENV GOPATH=/go
WORKDIR /go/src/peacemakr

RUN go install peacemakr-cli/main.go

FROM alpine

RUN apk update && apk upgrade && apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /go/bin/
COPY --from=builder /go/bin/main /go/bin/main

CMD ./main
