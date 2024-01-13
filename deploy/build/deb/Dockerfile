FROM debian:bookworm-slim
WORKDIR /root
RUN apt-get update \
    && apt-get -y install --no-install-recommends build-essential fakeroot devscripts \
    && groupadd --gid 1000 builder \
    && useradd --gid 1000 --uid 1000 --create-home builder \
    && apt-get clean \
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/* 
WORKDIR /home/builder
USER builder
COPY . .
