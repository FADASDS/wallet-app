FROM golang:alpine AS build-stage
ENV APP_HOME=/wallet-app
WORKDIR "$APP_HOME"
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/app

FROM alpine
WORKDIR /wallet-app
COPY --from=build-stage /wallet-app/app /wallet-app/app
ENV APP_PORT=4000
EXPOSE "$APP_PORT"
CMD ["/wallet-app/app"]

