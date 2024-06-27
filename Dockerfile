FROM golang AS builder
ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o echo

FROM scratch
COPY --from=builder /app/echo /echo
CMD ["./echo"]
