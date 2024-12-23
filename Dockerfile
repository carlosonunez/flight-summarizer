FROM scratch
ARG ARCH=arm64
COPY "out/flight_summarizer-linux-$ARCH" /app

ENTRYPOINT [ "/app" ]
