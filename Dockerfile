FROM golang:1.20

WORKDIR /service

ADD . .

RUN go get image-resizer
RUN apt-get update
RUN apt-get install -y libvips libvips-dev

RUN go build -o ./out/bin ./cmd/app/
WORKDIR /service/out

CMD ["./bin"]
