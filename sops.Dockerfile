FROM alpine:3.21
RUN apk add sops gpg gpg-agent
ENTRYPOINT [ "sops" ]
