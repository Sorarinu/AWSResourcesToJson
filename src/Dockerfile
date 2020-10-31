FROM golang:1.14.4-alpine

RUN apk update && \
    apk add --no-cache git make gcc musl-dev && \
    rm -rf /var/cache/apk/*

RUN mkdir -p /app
WORKDIR /app
COPY . .

RUN make

CMD ./main.handle