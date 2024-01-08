FROM golang:1.21

RUN go install github.com/cespare/reflex@latest

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download -x

COPY configs/reflex.conf /

ENTRYPOINT ["reflex", "-c", "/reflex.conf"]