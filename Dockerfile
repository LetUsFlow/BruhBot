FROM golang:1.18-bullseye AS build

WORKDIR /src/bruhbot

COPY . .

RUN go build -ldflags="-s -w" -o bin/bruhbot


FROM debian:bullseye-slim AS bin

WORKDIR /opt/bruhbot

COPY --from=build /src/bruhbot/bin/bruhbot ./
COPY --from=build /src/bruhbot/sounds/ ./sounds/

RUN apt-get update
RUN apt-get install ca-certificates ffmpeg --no-install-recommends -y

CMD [ "./bruhbot" ]
