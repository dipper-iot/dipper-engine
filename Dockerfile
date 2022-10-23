FROM alpine:3.12.1 as builder

COPY --from=golang:1.19-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"
RUN apk --no-cache add make git gcc libtool musl-dev

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /
RUN make build; rm -rf $GOPATH/pkg/mod


FROM alpine:3.12.1
RUN apk add --no-progress --no-cache ca-certificates
COPY --from=builder /dipper-engine /dipper-engine
ADD config.json /config.json
ENTRYPOINT [ "/dipper-engine" ]