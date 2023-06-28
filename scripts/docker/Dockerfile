FROM golang:alpine AS builder
LABEL maintainer="Niraj Georgian <ng@datawyse.io>"

# define system module and update them
RUN apk update && apk add --no-cache git

# Create user for accessing datawyse service.
RUN adduser -D -g '' serviceuser
WORKDIR /auth-service

ENV GO111MODULE=on

# install dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o datawyse main.go

# use scratch (base for a docker image)
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd

# set working directory
WORKDIR /root/
COPY --from=builder /auth-service .

# Use an unprivileged user.
USER serviceuser
EXPOSE 3000
CMD ["./auth-service", "serve"]