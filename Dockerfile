FROM golang:1.19.4-alpine3.17 AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

COPY . .

WORKDIR /usr/src/app/cmd/basic-docker

RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -v -o /usr/local/bin/deploystack ./...

FROM gcr.io/distroless/base AS runtime

COPY --from=builder /usr/local/bin/deploystack /deploystack

CMD [ "/deploystack" ]