FROM golang:1.17 AS buildstage
COPY placeholder.go /go/src
ENV CGO_ENABLED=0
RUN go build -o bin/ src/placeholder.go

FROM scratch
COPY --from=buildstage /go/bin/placeholder .
CMD ["/placeholder", "-listen", "systemd"]
