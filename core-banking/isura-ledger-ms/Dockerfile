##################################
# STEP 1 build executable binary #
##################################
FROM golang:1.26.1-alpine as builder

# Instala certificados CA e cria usuário não-privilegiado
RUN apk --no-cache add ca-certificates && \
    adduser -D -g '' appuser \

WORKDIR /app
COPY . .

# Fetch dependencies.
RUN go mod download

# Build otimizado: -w -s remove símbolos de debug
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags musl -ldflags="-w -s" -o /go/bin/server ./cmd/server/main.go

##############################
# STEP 2 time zone           #
##############################
FROM alpine:latest AS time-zone
RUN apk --no-cache add tzdata zip
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -q -r -0 /zoneinfo.zip .

##############################
# STEP 3 build a small image #
##############################
FROM scratch

# Importante: Copiar os certificados para conexões seguras
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar usuário para não rodar como root
COPY --from=builder /etc/passwd /etc/passwd
USER appuser

# Copy our static executable
COPY --from=builder /go/bin/server /go/bin/server

# Copy config json
COPY --from=builder /app/config.json /config.json

# Copy the zoneinfo.zip file from the time-zone stage.
COPY --from=time-zone /zoneinfo.zip /
ENV ZONEINFO /zoneinfo.zip

#Expose ports grpc server and http server
EXPOSE 50051 8080

ENTRYPOINT ["/go/bin/server"]
