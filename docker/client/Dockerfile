FROM golang:1.16-alpine AS base

FROM base as deps
WORKDIR "/chickaree-db"
ADD *.mod *.sum ./
RUN go mod download

FROM deps AS build-env
ADD cmd ./cmd
ADD chickaree ./chickaree
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -X main.docker=true" -o client cmd/client-server/main.go

FROM scratch AS prod
WORKDIR /
ENV PORT 6379
EXPOSE 6379
COPY --from=build-env /chickaree-db/client /
CMD ["/client"]