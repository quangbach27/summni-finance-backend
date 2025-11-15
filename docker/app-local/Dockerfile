FROM golang:1.25

RUN go install github.com/cespare/reflex@v0.3.1 && \
    go install github.com/go-delve/delve/cmd/dlv@v1.7.3

COPY reflex.conf /
COPY start.sh /

ENTRYPOINT ["reflex", "-c", "/reflex.conf"]
