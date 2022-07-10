# Build Stage
FROM golang:1.18.3-alpine3.16 AS buildStage
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.16
WORKDIR /app
COPY --from=buildStage /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]
