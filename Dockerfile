FROM golang:1.21 as builder
RUN adduser -u 10001 scratchuser

FROM scratch
COPY skyline /
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER scratchuser
ENTRYPOINT ["/skyline", "serve"]
