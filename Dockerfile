FROM golang:1.16-alpine AS builder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache make git

WORKDIR /app/
COPY . .
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=gosum.io+ce6e7565+AY5qEHUk/qmHc5btzW45JVoENfazw8LielDsaI+lEbq6
RUN go build ./example/consumer
RUN go build ./example/producer


FROM alpine:3
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk --no-cache add tzdata
ENV TZ=Asia/Shanghai

WORKDIR /usr/bin/
COPY --from=builder /app/consumer .
COPY --from=builder /app/producer .

ENTRYPOINT [ "sh", "-c", "while true; do sleep 1; done" ]
CMD []
