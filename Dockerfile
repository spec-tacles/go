FROM golang:alpine

WORKDIR /usr/bin/spectacles
COPY go.mod go.sum ./

RUN apk update && apk add git
RUN go mod download
RUN apk del git
COPY . .
RUN go build -o spectacles
CMD ["./spectacles"]
