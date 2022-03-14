# First stage (builder) to compile an executable.
# Firstly copy only go.mod and go.sum to download modules and cache the layer.
# It also creates user to make records in /etc/passwd to copy it into scratch image.
FROM golang:1.17.6-alpine3.15 AS builder
RUN adduser -S noroot
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /binary ./cmd

# Second stage without compiler in image and not-root user.
# Copies executable from first stage and declares CMD to execute it.
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
USER noroot
COPY --from=builder /binary /
ENTRYPOINT ["/binary"]