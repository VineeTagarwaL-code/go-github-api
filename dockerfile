FROM golang:1.23.2-alpine AS build

WORKDIR /app
ENV GOPROXY="https://goproxy.io,direct"

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=build /app/app .

EXPOSE 3000

CMD ["./app"]