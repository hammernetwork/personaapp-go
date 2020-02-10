FROM registry.gitlab.com/persona_app_online/devops/golangci:latest as build

ARG PACKAGE_NAME=personaapp-go
ARG APP_NAME=personapp
ARG PROJECT_NAMESPACE=persona_app_online

WORKDIR ./src/gitlab.com/${PROJECT_NAMESPACE}/${PACKAGE_NAME}

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile ./main.go ./
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

RUN make build && \
    cp ./bin/${APP_NAME} /usr/local/bin/ && \
    rm -rf /go/src/gitlab.com

FROM registry.gitlab.com/persona_app_online/devops/base_image:latest

COPY --from=build /usr/local/bin/${APP_NAME} /usr/local/bin/${APP_NAME}

ENV BIND 0.0.0.0:8000

EXPOSE 8000

ENTRYPOINT ["personapp"]
