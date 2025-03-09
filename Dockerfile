FROM golang:1.23.6-bullseye AS builder
WORKDIR /app
RUN apt-get update -qq && \
	apt-get install --no-install-recommends -y build-essential pkg-config python-is-python3 upx

# install zig toolchain
RUN wget https://ziglang.org/download/0.13.0/zig-linux-x86_64-0.13.0.tar.xz && \
	tar -xf zig-linux-x86_64-0.13.0.tar.xz && \
	mv zig-linux-x86_64-0.13.0 /usr/local/zig && \
	rm zig-linux-x86_64-0.13.0.tar.xz && \
	ln -s /usr/local/zig/zig /usr/local/bin/zig && \
	zig version

RUN apt-get install -y --no-install-recommends ca-certificates

COPY go.mod go.sum ./
RUN go version
RUN go mod tidy
COPY . .
RUN make -f tiny-bundle.mk build

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin .

EXPOSE 8080
CMD [ "/app/app_prod" ]