FROM golang:1.19 as builder

ARG ARG_GOPROXY
ENV GOPROXY $ARG_GOPROXY

WORKDIR /home/app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN make build

FROM orvice/go-runtime

ENV PROJECT_NAME sox
ENV CONFIG_PATH /etc/sox.toml

COPY --from=builder /home/app/bin/${PROJECT_NAME} /app/bin/${PROJECT_NAME}

ENTRYPOINT exec /app/bin/${PROJECT_NAME} -c ${CONFIG_PATH}