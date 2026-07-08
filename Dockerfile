# --- builder stage
FROM golang:1.26 AS builder

# environment
WORKDIR /workspace
ARG TARGETOS
ARG TARGETARCH

# install & cache dependencies
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# inject source code
COPY . .

# build manager binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

# --- final stage
FROM gcr.io/distroless/static:nonroot

# environment
WORKDIR /

# inject binary from builder stage
COPY --from=builder /workspace/manager .

# switch to non-root user
USER 65532:65532

# let's go!
ENTRYPOINT ["/manager"]
