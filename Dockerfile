FROM golang:1.19-alpine as builder

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download
RUN CGOENABLED=0 GOOS="linux" go build -o service ./cmd/api/main.go


# Intermediate stage: Build the binary
FROM golang:1.19-alpine as runner


WORKDIR /app
ENV config=docker

COPY --from=builder /app/service /app/service
COPY db /app/db
#COPY ./config/config-local.yml /app/config/config-local.yml

EXPOSE 5000
EXPOSE 5555
EXPOSE 7070

ENTRYPOINT /app/service