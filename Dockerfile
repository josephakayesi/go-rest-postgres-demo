FROM golang:1.23.2-alpine3.20 AS builder

LABEL author="Joseph Akayesi <josephakayesi@gmail.com>"

WORKDIR /dist

COPY . .

RUN go mod download

RUN go mod verify

ENV CGO_CPPFLAGS="-D_FRtify_SOURCE=2 -fstack-protector-all"
ENV GOFLAGS="-buildmode=pie"
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/main ./cmd/main.go

FROM alpine:3.19
# FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=golang /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /etc/passwd /etc/passwd
COPY --from=golang /etc/group /etc/group

COPY --from=builder /bin/main /bin/main
# COPY --from=builder /dist/keys /keys

EXPOSE 5000

ENTRYPOINT ["/bin/main"]
