FROM golang:1.23.2-alpine3.20 AS build
WORKDIR /usr/app

RUN apk update && apk add --no-cache make build-base gcc sqlite-dev

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make tidy

ENV GOCACHE=/usr/app/.cache/go-build
RUN --mount=type=cache,target="/usr/app/.cache/go-build" go build -ldflags='-s -w' -o=./bin/tinyurl ./cmd/api

FROM alpine:3.20 AS final
ENV APP_HOME=/home/app
WORKDIR ${APP_HOME}

RUN apk add --no-cache sqlite-libs
COPY --from=build /usr/app/bin/tinyurl ./

RUN addgroup -S app && adduser -S app -G app
RUN chown -R app:app $APP_HOME

USER app

CMD [ "./tinyurl" ]