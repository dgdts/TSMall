FROM hubzhidc.seasungame.com/tc-devops/ci:golang-1.21.10-jdk-8-centos-7.9 as builder
LABEL stage=gobuilder

ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0

ARG ARTIFACT_ID

WORKDIR /build
COPY ./$ARTIFACT_ID/output .



FROM hubzhidc.seasungame.com/tc-devops/os_base:alpine-3.16

ARG ARTIFACT_ID
USER root 

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build /application/$ARTIFACT_ID

ENV TZ=Asia/Shanghai
WORKDIR /application/$ARTIFACT_ID

# COPY /usr/share/zoneinfo /usr/share/zoneinfo
# COPY ./$ARTIFACT_ID/output .
EXPOSE 8888
RUN chmod +x /application/$ARTIFACT_ID/bin/$ARTIFACT_ID && chmod 4755 /bin/busybox && apk update && apk add curl && apk add busybox-extras

ENTRYPOINT ["./bin/ts_mall_service"]