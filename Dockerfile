# Builder
FROM golang AS builder

# ENV GO111MODULE=off
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src
COPY ./http-server/ ./
RUN go get && go build -o http-server

# Target
FROM scratch

ENV VERSION=1.0
ENV BIND_PORT=80
ENV LOG_LEVEL=INFO
EXPOSE 80

WORKDIR /app
COPY --from=builder /go/src/http-server .
ENTRYPOINT ["/app/http-server"]