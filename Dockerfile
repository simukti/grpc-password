### BUILD STAGE
FROM golang:1.9 as builder
RUN apt-get update && apt-get install -y golang-glide build-essential
ENV APP_PATH    github.com/simukti/grpc-password
ENV CGO_ENABLED 1
ENV GOOS        linux
ENV GOARCH      amd64
ENV LIBSODIUM_VERSION 1.0.15
ENV LIBSODIUM_DIR     /tmp/sodium
ENV LIBSODIUM_SERVER  https://github.com/jedisct1/libsodium/releases/download
RUN mkdir -p ${LIBSODIUM_DIR} && \
    cd ${LIBSODIUM_DIR} && \
    curl -L -o libsodium-${LIBSODIUM_VERSION}.tar.gz ${LIBSODIUM_SERVER}/${LIBSODIUM_VERSION}/libsodium-${LIBSODIUM_VERSION}.tar.gz && \
    tar zxfv libsodium-${LIBSODIUM_VERSION}.tar.gz && \
    cd libsodium-${LIBSODIUM_VERSION}/ && \
    ./configure && \
    make && make install && \
    mv ./src/libsodium /usr/local/ && \
    cd / && \
    rm -Rfv ${LIBSODIUM_DIR}
WORKDIR $GOPATH/src/$APP_PATH
ADD . $GOPATH/src/$APP_PATH
RUN glide --version
RUN glide install
### warning on static-linked app is intended, see: https://github.com/golang/go/issues/21421
RUN go build -a --ldflags '-s -w -linkmode external -extldflags "-static"' -o /go/bin/grpc-password
### uncomment these two lines if you want more compact app binary size
### ps. BEWARE: it will took about 2-4 minutes on docker build and consume more memory on runtime
#RUN apt-get -y install upx-ucl
#RUN /usr/bin/upx-ucl -v --brute /go/bin/credential

### FINAL STAGE
FROM alpine:3.6
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/grpc-password ./grpc-password
CMD ["./grpc-password"]