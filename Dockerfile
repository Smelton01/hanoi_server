FROM golang:1.16-alpine AS builder
RUN mkdir /build 
ADD go.mod go.sum main.go /build/
WORKDIR /build
RUN go build -o hanoi

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/hanoi /app/
COPY views/ /app/views
WORKDIR /app
CMD ["./hanoi"] 
