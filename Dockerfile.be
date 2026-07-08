FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder
ARG TARGETARCH
WORKDIR /app
COPY ors-be/go.mod ors-be/go.sum ./ors-be/
RUN cd ors-be && go mod download
COPY ors-be/ ./ors-be/
RUN cd ors-be && CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o bin/ors-be ./cmd/server

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/ors-be/bin/ors-be /usr/local/bin/ors-be
EXPOSE 8080
CMD ["ors-be"]
