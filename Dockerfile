ARG DEP_BASE_VERSION
ARG CONTAINER_REPO
FROM ${CONTAINER_REPO}/alpine:${DEP_BASE_VERSION}

EXPOSE 8080

RUN apk add --no-cache --upgrade tzdata

COPY res/fixtures ./res/fixtures
COPY build/linux/golang-api-service .

CMD ["/golang-api-service"]
