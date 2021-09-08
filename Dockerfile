FROM golang:1.17 AS builder
WORKDIR /src/github.com/foroozf001/gotcp
COPY /src/github.com/foroozf001/gotcp .
RUN go get github.com/foroozf001/gotcp && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /src/github.com/foroozf001/gotcp/app .
CMD ["./app"]