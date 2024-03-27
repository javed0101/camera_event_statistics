FROM golang:1.21.4-alpine3.18 AS deps

RUN apk add --no-cache bash
RUN apk add --no-cache shadow
RUN apk add --no-cache git make

ENV HOME /usr/src/app
WORKDIR $HOME

COPY --from=harbor.eencloud.com/vms/goeen:1.0.109 /usr/src/app/go/src/github.com/eencloud/goeen /usr/src/goeen
RUN cd /usr/src/goeen && \
    go mod download

COPY ./go.mod $HOME/
COPY ./go.sum $HOME/
RUN cd $HOME/ && go mod download

CMD [ "sleep", "infinity"]

FROM deps AS build

ENV HOME /usr/src/app
WORKDIR $HOME

COPY . $HOME/
COPY . $HOME
ENV GOFLAGS=-tags=no_dhash,no_mpack
RUN cd $HOME && make build

CMD [ "sleep", "infinity"]

FROM alpine:3.18.2 AS prod

RUN apk add --no-cache bash
RUN apk add --no-cache shadow
RUN apk add --no-cache vim

ENV HOME /usr/src/app
WORKDIR $HOME

RUN mkdir -p "$HOME/config"

COPY --from=build $HOME/build/cameraevent $HOME
COPY --from=build $HOME/run.sh $HOME/run.sh
COPY --from=build $HOME/config/config.json $HOME/config

CMD ["sh", "run.sh"]