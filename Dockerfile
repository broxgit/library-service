FROM golang:alpine as builder

RUN apk update && apk add git

RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/libraryservice/main.go

FROM golang

WORKDIR /app
COPY --from=builder build/main .

ENTRYPOINT [ "./main" ]