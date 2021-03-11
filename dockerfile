FROM golang:latest AS go-builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=go-builder /app/static/ ./static/
COPY --from=go-builder /app/views/ ./views/
COPY --from=go-builder /app/main ./app

CMD [ "./app/main" ]