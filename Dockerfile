FROM golang:1.16-alpine AS build-env

WORKDIR /src

COPY . .
RUN go build -o cmd/main cmd/main.go


FROM alpine:edge

WORKDIR /root

COPY --from=build-env /src/cmd/main /root/main

RUN chmod +x /root/main

EXPOSE 8888

CMD ["/root/main"]
