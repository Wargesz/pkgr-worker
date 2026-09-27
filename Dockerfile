FROM ghcr.io/void-linux/void-glibc-full:20260901r1

WORKDIR /opt/pkgr/

ENV HOST=
RUN xbps-install -Syu wget tar curl zip file make gcc go
RUN mkdir build
COPY . .

CMD ["go", "run", "."]
