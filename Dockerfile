FROM golang AS builder

ENV GO111MODULE=off \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

WORKDIR /src
COPY . .
RUN go build -o myserver .

FROM scratch
WORKDIR /app
COPY --from=builder /src/myserver .
EXPOSE 80
ENTRYPOINT ["/app/myserver"]