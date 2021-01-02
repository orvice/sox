FROM golang:1.15 as builder

ARG ARG_GOPROXY
ENV GOPROXY $ARG_GOPROXY

WORKDIR /home/app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN make build

FROM orvice/go-runtime

ENV PROJECT_NAME twitter2mastodon
ENV CONFIG_PATH /etc/twitter2mastodon.toml

COPY --from=builder /home/app/bin/${PROJECT_NAME} /app/bin/${PROJECT_NAME}

ENTRYPOINT exec /app/bin/${PROJECT_NAME} -c ${CONFIG_PATH}