FROM registry.opensuse.org/opensuse/tumbleweed:latest AS build-stage
RUN zypper install --no-recommends --auto-agree-with-product-licenses -y git go make
COPY . daikin-gomod
RUN cd daikin-gomod && make update && make tidy && make

FROM registry.opensuse.org/opensuse/busybox:latest
LABEL maintainer="Thorsten Kukuk <kukuk@thkukuk.de>"

ARG BUILDTIME=
ARG VERSION=unreleased
LABEL org.opencontainers.image.title="Exports Daikin AC values as metrics for Prometheus"
LABEL org.opencontainers.image.description="Exports Daikin AirConditioner values as metrics for Prometheus"
LABEL org.opencontainers.image.created=$BUILDTIME
LABEL org.opencontainers.image.version=$VERSION

COPY --from=build-stage /daikin-gomod/bin/daikin-ac-exporter /usr/local/bin
COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/usr/local/bin/daikin-ac-exporter"]
