ARG base_app_dir_default=/go/src/github.com/sandbreaker/goservice
ARG app_name_default=service-api

FROM golang:1.13.7 as builder

ARG base_app_dir_default
ARG app_name_default
ENV BASE_APP_DIR=$base_app_dir_default
ENV APP_NAME=$app_name_default
RUN echo ${BASE_APP_DIR}
RUN echo ${APP_NAME}

RUN apt-get update && apt-get install -y --no-install-recommends \
        g++ \
        gcc \
        libc6-dev \
        ca-certificates \       
        curl \
        git \
        make \
        cmake \
        build-essential \
        && rm -rf /var/lib/apt/lists/*

# compile static binary
RUN mkdir -p ${BASE_APP_DIR}
WORKDIR ${BASE_APP_DIR}
COPY . .
RUN make build-static

FROM golang:1.13.7-alpine 
WORKDIR /app

ARG base_app_dir_default
ARG app_name_default
ENV BASE_APP_DIR=$base_app_dir_default
ENV APP_NAME=$app_name_default
RUN echo ${BASE_APP_DIR}
RUN echo ${APP_NAME}

COPY --from=builder ${BASE_APP_DIR}/bin/linux_amd64/${APP_NAME} /app
COPY --from=builder ${BASE_APP_DIR}/.cfg /app/.cfg

# Run Binary
RUN pwd

EXPOSE 8088
ENTRYPOINT ./${APP_NAME}
