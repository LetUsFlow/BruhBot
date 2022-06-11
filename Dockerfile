FROM golang:1.18-bullseye AS build

COPY . /app
WORKDIR /app

RUN go build -ldflags="-s -w" -o bin/bruhbot

FROM alpine:3

WORKDIR /app

COPY --from=build /app/bin/bruhbot /app/bruhbot
COPY --from=build /app/sounds/ /app/sounds/

RUN apk add --no-cache ffmpeg gcompat && \
    rm -fr /var/cache/apk/*

USER 1000

CMD [ "/app/bruhbot" ]
