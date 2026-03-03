FROM golang:1.25-bookworm AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/sshconfctl ./cmd/sshconfctl

FROM gcr.io/distroless/static-debian12
COPY --from=build /bin/sshconfctl /usr/local/bin/sshconfctl
ENTRYPOINT ["/usr/local/bin/sshconfctl"]
