FROM golang:1.17-alpine AS builder
RUN mkdir /build
RUN apk --update --no-cache add g++
ADD go.mod go.sum main.go  /build/
COPY grokify/ /build/grokify
WORKDIR /build
RUN apk update && \
    apk add nano
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/docker /app/
COPY views/ /app/views
WORKDIR /app
CMD ["./docker"]

 



