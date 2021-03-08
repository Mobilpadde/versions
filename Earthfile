build:
    FROM golang:alpine
    WORKDIR /

    RUN apk add g++

    COPY . .
    RUN go mod download
    RUN go build -o build/versions .

    SAVE ARTIFACT build/versions /versions

versions:
    FROM node:lts-alpine
    WORKDIR /

    RUN apk add chromium git

    RUN npm i -g yarn --force
    COPY fonts ./fonts/
    COPY +build/versions .

    VOLUME ["/repo", "/out", "/screendumps"]

    ENTRYPOINT ["./versions", "-repo", "/repo", "-out", "/out", "-dump", "/screendumps"]
    CMD ["./versions", "-repo", "/repo", "-out", "/out", "-dump", "/screendumps"]

    SAVE IMAGE versions