FROM golang:alpine AS builder

WORKDIR /app
COPY ./src .

RUN go get
RUN go build -o do-dyndns

FROM alpine
WORKDIR /
COPY --from=builder /app/do-dyndns /do-dyndns
CMD ["/do-dyndns"]
