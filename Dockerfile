FROM golang:1.21

WORKDIR /build
COPY . ./
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -o cf_ddns ./

FROM alpine:3.12
COPY --from=0 /build/cf_ddns /bin/

ENTRYPOINT [ "/bin/cf_ddns" ]