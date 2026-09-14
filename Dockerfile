FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./ 
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /kusuriya ./cmd/kusuriya

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /kusuriya /kusuriya
EXPOSE 8080
ENTRYPOINT ["/kusuriya"]