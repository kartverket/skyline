FROM golang:1.23

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -tags timetzdata

ENTRYPOINT ["/app/skyline", "serve"]