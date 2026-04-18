# stage 1: build
FROM golang:1.25 AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o jade .
# install goose binary
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest
# stage 2: final image
FROM alpine:latest 

RUN apk add --no-cache ca-certificates

WORKDIR /app 

# copy binary file 
COPY --from=builder /app/jade .
#copy goose
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# copy static files
COPY --from=builder /app/frontend/ ./frontend

# copy migrations 
COPY --from=builder /app/sql/schema ./sql/schema

# copy entrypoint
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT [ "./entrypoint.sh" ]