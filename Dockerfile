FROM golang:alpine as builder
LABEL stage=builder

WORKDIR /
RUN apk add g++

COPY . .
RUN go mod download
RUN go build -o versions .

# FROM zenika/alpine-chrome:with-node

# WORKDIR /tmp
# RUN mkdir npm
# RUN npm config set prefix /tmp/npm
# RUN npm i -g yarn

FROM node:lts-alpine

WORKDIR /
RUN apk add chromium git

RUN npm i -g yarn --force
# COPY --from=builder --chown=chrome /versions .
COPY --from=builder /versions .
COPY fonts ./fonts/

VOLUME ["/repo", "/out", "/screendumps"]
ENTRYPOINT ["./versions", "-repo", "/repo", "-out", "/out", "-dump", "/screendumps"]