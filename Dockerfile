FROM golang as compiler
RUN CGO_ENABLED=1 go get -a -ldflags '-s' \
github.com/YKatrechko/smdbot

FROM progrium/busybox
RUN opkg-install ca-certificates
RUN mkdir -p /bot
WORKDIR /bot
VOLUME /bot
COPY --from=compiler /go/bin/smdbot ./
COPY .config.json .
CMD ["./smdbot"]
