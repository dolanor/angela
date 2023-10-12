# syntax=docker/dockerfile:1

FROM golang:1.21 AS build

WORKDIR /src

ENV GOMODCACHE /root/.cache/gocache
RUN --mount=target=. --mount=type=cache,target=/root/.cache \
    CGO_ENABLED=0 go build -o /bin/filed ./cmd/filed


FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /bin/filed /filed

CMD [ "/filed" ]
