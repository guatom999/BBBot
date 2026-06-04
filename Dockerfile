FROM golang:1.24-bullseye AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/app

# FROM gcr.io/distroless/static-debian11
FROM debian:bullseye-slim 
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /bin/app /bin

CMD ["/bin/app"]