FROM golang:1.19.3 as builder

ENV APP_HOME /go/src/questapp

WORKDIR "$APP_HOME"
COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o questapp

FROM golang:1.19.3

ENV APP_HOME /go/src/questapp
ENV REST_PORT 8010

RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=builder "$APP_HOME"/questapp $APP_HOME

EXPOSE 8010
CMD ["./questapp"]