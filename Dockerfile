FROM golang:1.14.6 AS build
WORKDIR /build
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o codebuild

FROM alpine:3.12.0 as certs
RUN apk --update add ca-certificates

FROM scratch
WORKDIR /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /build/codebuild /codebuild
ENTRYPOINT ["./codebuild"]
