# syntax=docker/dockerfile:1
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
  libxkbcommon0 \
  ca-certificates \
  git \
  golang \
  unzip \
  libc++1 \
  vim \
  && apt-get clean autoclean

# Ensure UTF-8 encoding
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

WORKDIR /workspace

COPY . /workspace
