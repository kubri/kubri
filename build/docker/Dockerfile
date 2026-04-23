FROM gcr.io/distroless/static:latest

ARG TARGETPLATFORM

COPY $TARGETPLATFORM/kubri /usr/bin/

ENV KUBRI_PATH=/kubri

ENTRYPOINT ["kubri"]
