FROM aviasales/build-go:1.12-alpine as build

WORKDIR $GOPATH/src/github.com/KosyanMedia/burlesque/exporter
COPY . .
RUN dep check || dep ensure --vendor-only -v
RUN go build -o /go/bin/burlesque_exporter .

FROM alpine:3.9
MAINTAINER  Maxim Pogozhiy <foxdalas@gmail.com>

RUN apk --no-cache add ca-certificates tzdata
COPY --from=build /go/bin/burlesque_exporter /bin/

ENTRYPOINT ["/bin/burlesque_exporter"]
EXPOSE     9118
