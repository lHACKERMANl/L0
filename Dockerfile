FROM golang:1.21.0 AS build-env

RUN go install github.com/go-delve/delve/cmd/dlv@v1.21.0

COPY config.yaml config.yaml
ADD . /dockerdev
WORKDIR /dockerdev

RUN go build -gcflags="all=-N -l" -o /server


FROM debian:latest

EXPOSE 8000 4000

WORKDIR /
COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /server /
COPY --from=build-env /dockerdev /
#COPY --from=build-env /dockerdev/static /static

#CMD ["/dlv", "--listen=:4000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server"]
CMD ["/server"]