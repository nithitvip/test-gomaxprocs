FROM golang:1.20-alpine AS build

WORKDIR /test-app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o /go/bin/test_app .

FROM alpine:3.18
COPY --from=build /go/bin/test_app /go/bin/test_app
ENTRYPOINT ["/go/bin/test_app"]
