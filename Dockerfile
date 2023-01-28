FROM golang:1.19.0-buster AS builder

WORKDIR /gloader

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o gloader ./cmd/gloader

FROM debian:buster-slim

COPY --from=builder /gloader/gloader ./gloader

ENTRYPOINT [ "./gloader" ]