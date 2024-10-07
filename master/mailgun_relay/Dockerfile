FROM golang:1.22.5-alpine3.20 AS go-builder

ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

WORKDIR /build
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -o ./bin/mail-relay ./main.go


FROM alpine:3.20 AS runner
WORKDIR /app/bin
COPY --from=go-builder /build/bin/* /app/bin/
COPY --from=go-builder /build/template/* /app//bin/template/
COPY --from=go-builder /build/bawana.png /app/bin/
COPY --from=go-builder /build/receivers.json /app/bin/
COPY --from=go-builder /build/webhook_req.json /app/bin/
COPY --from=go-builder /build/test-run.sh /app/bin/
COPY --from=go-builder /build/.env /app/bin/

CMD ["/app/bin/mail-relay"]
