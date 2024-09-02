ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine as builder

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  small-user

WORKDIR $GOPATH/src/smallest-golang/app/
COPY . .
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /run-app .

FROM gcr.io/distroless/static-debian12

COPY --from=builder /run-app .
CMD ["./run-app"]
