### Telegram WEBM-Converter DockerFile ###

FROM golang:1.22.1-alpine
RUN mkdir /telegram-webm-converter && apk add --no-cache ffmpeg
ADD . /telegram-webm-converter
WORKDIR /telegram-webm-converter
RUN go build -o telegram-webm-converter .
LABEL Name=telegram-webm-converter Version=0.0.1
COPY config.toml /config/config.toml
COPY media /media
CMD ["./telegram-webm-converter"]