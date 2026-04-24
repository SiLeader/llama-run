FROM golang:1.26-alpine AS builder

LABEL authors="cerussite"

ENV CGO_ENABLED=0

WORKDIR /work

COPY . .

RUN go build -o /llama-run

FROM scratch

COPY --from=builder /llama-run /usr/local/bin/llama-run

EXPOSE 8080

CMD ["/usr/local/bin/llama-run"]
