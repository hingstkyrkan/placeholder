FROM docker.io/golang:1.20 AS buildstage
COPY placeholder.go /go/src
ENV CGO_ENABLED=0
RUN go build -ldflags="-w -s" -o bin/ src/placeholder.go

FROM scratch
COPY --from=buildstage /go/bin/placeholder .
EXPOSE 8080
CMD ["/placeholder", "-listen", "tcp@:8080"]
