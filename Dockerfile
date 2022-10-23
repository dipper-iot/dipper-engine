FROM golang:1.19-alpine as builder

ENV GO111MODULE=on
ENV GOFLAGS=" -ldflags '-w'"
ENV GOPROXY=direct
ENV GOSUMDB=off

COPY . .

RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build

FROM alpine:3.14 as run
RUN apk add --no-progress --no-cache ca-certificates
COPY --from=builder dipper-engine /dipper-engine
ADD config.json /config.json
ENTRYPOINT [ "/dipper-engine" ]