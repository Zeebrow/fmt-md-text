FROM golang:1.17

WORKDIR /app
COPY . .
RUN mkdir /output && go build -o /output/fmt-md-text .
