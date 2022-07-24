FROM golang AS builder

ENV GOOS=linux \
	GOARCH=amd64

WORKDIR /src
COPY ./http-server ./
RUN go get && go build

FROM scratch
WORKDIR /app
COPY --from=builder /src/http-server .
EXPOSE 80
ENTRYPOINT ["/app/http-server"]