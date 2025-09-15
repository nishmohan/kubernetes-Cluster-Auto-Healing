FROM golang:1.24.5 AS builder

WORKDIR /src

#1) copy only mod files first to leverage Docker layer cache
COPY go.mod go.sum ./
RUN go mod download

#2)now copy the rest of the source
COPY . .

#3) build
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /autoheal ./main.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /autoheal /autoheal
USER 65532:65532
ENTRYPOINT ["/autoheal"]
