FROM golang:1.17.6-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY records/* records/
COPY *.go ./

RUN go build -o /ha-postgres

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /ha-postgres /ha-postgres

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/ha-postgres"]