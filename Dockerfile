FROM golang:1.14.6 AS build
WORKDIR /build
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o codebuild

FROM alpine:3.12.0
RUN apk add --update ca-certificates
COPY --from=build /build/codebuild /bin/codebuild
ENTRYPOINT ["codebuild"]
