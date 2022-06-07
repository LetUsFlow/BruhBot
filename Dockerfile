FROM golang:1.18-bullseye AS build

COPY . /app
WORKDIR /app

RUN go build -ldflags="-s -w" -o bin/bruhbot

FROM alpine as compressor

COPY --from=build /app/bin/bruhbot /app/bin/bruhbot

RUN apk add upx && \
    upx --lzma --best /app/bin/bruhbot

FROM frolvlad/alpine-glibc

WORKDIR /app

COPY --from=compressor /app/bin/bruhbot /app/bruhbot
COPY --from=build /app/sounds/ /app/sounds/

RUN apk add --no-cache ffmpeg && \
    rm -fr /var/cache/apk/*

USER 1000

CMD [ "/app/bruhbot" ]
