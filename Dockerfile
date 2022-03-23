FROM golang:alpine as builder
RUN apk add git
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go build -o app
FROM alpine
COPY --from=builder /build/app /app/
WORKDIR /app
EXPOSE 8000
CMD ["./app"]