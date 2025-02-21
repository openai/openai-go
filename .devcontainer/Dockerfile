# syntax=docker/dockerfile:1
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
  libxkbcommon0 \
  ca-certificates \
  git \
  unzip \
  libc++1 \
  vim \
  curl \
  procps \
  && apt-get clean autoclean

RUN curl -OL https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && \
  tar -C /usr/local -xzvf go1.24.0.linux-amd64.tar.gz && \
  rm go1.24.0.linux-amd64.tar.gz
ENV PATH="$PATH:/usr/local/go/bin"

# Ensure UTF-8 encoding
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8
