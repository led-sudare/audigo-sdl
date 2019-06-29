#build stage
FROM golang:1.12-alpine3.10 AS builder
RUN apk add --no-cache git \
    gcc libc-dev libx11-dev \
    sdl2-dev sdl2_mixer-dev
WORKDIR /app
COPY . /app
ENV GO111MODULE=on
RUN go build -a

#final stage
FROM alpine:3.10
RUN apk --no-cache add ca-certificates \
    sdl2 sdl2_mixer
ARG port="80"
ENV _PORT ${port}
WORKDIR /app
COPY --from=builder /app/asset /app/asset
COPY --from=builder /app/audigo-sdl /app/
CMD /app/audigo-sdl ${_PORT}
