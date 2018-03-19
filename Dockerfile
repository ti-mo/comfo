FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY comfo_linux_amd64 .
CMD ["./comfo_linux_amd64"]
