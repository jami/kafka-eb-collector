FROM golang:1.13.9-alpine3.10 as builder
ARG LIBRDKAFKA_VERSION=v1.3.0
ARG LIBCYRUSSASL_VERSION=2.1.27

RUN apk add --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
    zlib-dev \
    musl-dev \
    make \
    git \
    bash \
    g++ \
    go \
    zstd-static zstd-libs zstd-dev \
    openssl-dev \
    build-base \
    curl \
    ca-certificates \
    librdkafka-dev \
    lz4-dev 

# Build Cyrus SASL from source
RUN cd $(mktemp -d) \
    && curl -sL "https://github.com/cyrusimap/cyrus-sasl/releases/download/cyrus-sasl-$LIBCYRUSSASL_VERSION/cyrus-sasl-$LIBCYRUSSASL_VERSION.tar.gz" | \
    tar -xz --strip-components=1 -f - \
    && ./configure \
        --prefix=/usr --disable-sample --disable-obsolete_cram_attr --disable-obsolete_digest_attr --enable-static --disable-shared --disable-checkapop --disable-cram --disable-digest --enable-scram --disable-otp --disable-gssapi --with-dblib=none --with-pic \
    && make \
    && make install

# Build librdkafka from source
#RUN cd $(mktemp -d) \
#    && curl -sL "https://github.com/edenhill/librdkafka/archive/$LIBRDKAFKA_VERSION.tar.gz" | \
#    tar -xz --strip-components=1 -f - \
#    && ./configure --disable-sasl \
#    && make -j \
#    && make install

COPY . /go/src/github.com/jami/kafka-eb-collector
WORKDIR /go/src/github.com/jami/kafka-eb-collector

RUN export GOPATH=/go && make deps && make build/linux

# runtime container
FROM alpine:3.10

ENTRYPOINT []
WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/jami/kafka-eb-collector/bin/kafka-collector /bin/kafka-collector

CMD [ "/bin/kafka-collector" ]