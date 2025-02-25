FROM golang:alpine AS build-stage
ENV APP_HOME=/wallet-app
WORKDIR "$APP_HOME"
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go test -v ./tests/
RUN go build -o app ./cmd/app

FROM alpine
WORKDIR /wallet-app
COPY --from=build-stage /wallet-app/app /wallet-app/app
EXPOSE 8080
CMD ["/wallet-app/app"]

