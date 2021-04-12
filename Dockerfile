FROM golang:1.16beta1-alpine3.12 AS build
# Support CGO and SSL
RUN apk --no-cache add gcc g++ make
RUN apk add git
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/social ./cmd/social/main.go
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/dialog ./cmd/dialog/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
COPY --from=build /go/src/app/web web
EXPOSE 3000
#ENTRYPOINT /go/bin/social