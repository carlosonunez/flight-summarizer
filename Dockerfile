FROM --platform=$BUILDPLATFORM curlimages/curl:8.11.1 AS timezone_data
COPY ./include/retrieve_tz.sh /script.sh
USER root
RUN chmod 755 /script.sh
RUN /script.sh

FROM scratch
ARG PROJECT_NAME TARGETOS TARGETARCH
COPY --from=timezone_data /data /data
COPY --from=alpine:3.21 /etc/ssl/cert.pem /etc/ssl/cert.pem
COPY "out/$PROJECT_NAME-$TARGETOS-$TARGETARCH" /app
ENTRYPOINT [ "/app" ]
