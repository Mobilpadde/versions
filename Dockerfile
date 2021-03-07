FROM golang:alpine as builder
LABEL stage=builder

WORKDIR /
RUN apk add g++

COPY . .
RUN go mod download
RUN go build -o versions .

FROM node:lts-alpine

WORKDIR /
RUN apk add chromium git

RUN npm i -g yarn --force
COPY --from=builder /versions .
COPY fonts ./fonts/

VOLUME ["/repo", "/out", "/screendumps"]

ENTRYPOINT ["./versions", "-repo", "/repo", "-out", "/out", "-dump", "/screendumps"]
CMD ["./versions", "-repo", "/repo", "-out", "/out", "-dump", "/screendumps"]