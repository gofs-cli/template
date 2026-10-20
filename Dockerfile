FROM golang:1.24 AS build

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/server

FROM gcr.io/distroless/static-debian12:latest

COPY --from=build /go/bin/app /
EXPOSE 8080
CMD ["/app"]
