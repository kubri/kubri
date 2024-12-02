FROM alpine:latest

RUN apk update && \
  apk add --no-cache vsftpd && \
  mkdir -p /var/log/vsftpd /home/ftpuser && \
  touch /var/log/vsftpd/xferlog && \
  touch /etc/vsftpd.userlist && \
  mkdir -p /home/ftpuser

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENV FTP_USER=ftpuser \
  FTP_PASS=password \
  PASV_ADDRESS=0.0.0.0 \
  PASV_MIN_PORT=21100 \
  PASV_MAX_PORT=21110 \
  FTP_PORT=21

EXPOSE 21 21100-21110

ENTRYPOINT ["/entrypoint.sh"]
