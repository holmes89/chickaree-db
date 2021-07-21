FROM golang:1.16-alpine AS base

FROM base as deps
WORKDIR "/chickaree-db"
ADD *.mod *.sum ./
RUN go mod download

FROM deps AS build-env
ADD cmd ./cmd
ADD chickaree ./chickaree
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -X main.docker=true" -o storage cmd/storage-server/main.go
CMD ["/chickaree-db/storage"]

FROM scratch AS prod

WORKDIR /
EXPOSE 8040
EXPOSE 8041
COPY --from=build-env /chickaree-db/storage /
CMD ["/storage"]