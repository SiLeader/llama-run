FROM golang:1.26-alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /work

COPY . .

RUN go build -o /llama-run

FROM ubuntu:24.04 AS llama-cpp

RUN apt-get update && \
    apt-get install -y ca-certificates curl tar && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /work

ENV LLAMA_CPP_TAG=b8864
ENV LLAMA_CPP_CHECKSUM=06f402cdf89bd0622916c5813505657740dc4290ba7c37c176f1c1c2ef48504e

RUN echo "Downloading llama.cpp version ${LLAMA_CPP_TAG}..." && \
    echo "Expected SHA256 checksum: ${LLAMA_CPP_CHECKSUM}" && \
    echo "Downloading from:  https://github.com/ggml-org/llama.cpp/releases/download/${LLAMA_CPP_TAG}/llama-${LLAMA_CPP_TAG}-bin-ubuntu-x64.tar.gz" && \
    curl -Lo llama.tar.gz "https://github.com/ggml-org/llama.cpp/releases/download/${LLAMA_CPP_TAG}/llama-${LLAMA_CPP_TAG}-bin-ubuntu-x64.tar.gz" && \
    echo "${LLAMA_CPP_CHECKSUM}  llama.tar.gz" > checksum.txt && \
    cat checksum.txt && \
    sha256sum -c checksum.txt && \
    tar -xzf llama.tar.gz --strip-components=1 && \
    rm llama.tar.gz

FROM ubuntu:24.04

ARG git_commit="unknown"

LABEL org.opencontainers.image.source="https://github.com/sileader/llama-run.git" \
      org.opencontainers.image.url="https://github.com/sileader/llama-run" \
      org.opencontainers.image.revision="${git_commit}" \
      org.opencontainers.image.licenses="MIT"

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /llama-run /usr/local/bin/llama-run
COPY --from=llama-cpp /work /opt/llama.cpp
COPY config.yaml /etc/llama-run/config.yaml

WORKDIR /opt/llama.cpp

EXPOSE 8080
USER 1000:1000

CMD ["/usr/local/bin/llama-run"]
