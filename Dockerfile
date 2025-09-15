# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/polyschedule ./cmd/polyschedule-backend

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /
USER nonroot:nonroot
COPY --from=builder /bin/polyschedule /polyschedule
EXPOSE 8080
ENTRYPOINT ["/polyschedule"]


