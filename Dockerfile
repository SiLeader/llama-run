FROM golang:1.26-alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /work

COPY . .

RUN go build -o /llama-run

FROM ubuntu:24.04 AS llama-cpp

ARG llama_cpp_tag="b8925"
ARG llama_cpp_checksum="e50407e42b1db107e7a6781efc6b9a1b37c3958931b3d2d0e6bafd3ca6b8c62a"
ARG llama_cpp_arch="x64"

RUN apt-get update && \
    apt-get install -y ca-certificates curl tar && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /work

RUN echo "Downloading llama.cpp version ${llama_cpp_tag}..." && \
    echo "Expected SHA256 checksum: ${llama_cpp_checksum}" && \
    echo "Downloading from:  https://github.com/ggml-org/llama.cpp/releases/download/${llama_cpp_tag}/llama-${llama_cpp_tag}-bin-ubuntu-${llama_cpp_arch}.tar.gz" && \
    curl -Lo llama.tar.gz "https://github.com/ggml-org/llama.cpp/releases/download/${llama_cpp_tag}/llama-${llama_cpp_tag}-bin-ubuntu-${llama_cpp_arch}.tar.gz" && \
    echo "${llama_cpp_checksum}  llama.tar.gz" > checksum.txt && \
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
COPY docker/config.yaml /etc/llama-run/config.yaml

WORKDIR /opt/llama.cpp

EXPOSE 8080

CMD ["/usr/local/bin/llama-run"]
