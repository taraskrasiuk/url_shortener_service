FROM golang:1.24-bullseye AS build
## create a non-root user to run the application
# RUN groupadd -g 1000 appgroup \
#     && useradd -u 1000 -g appgroup -s /bin/bash -m appuser
# set a working directory
WORKDIR /app
# copy required go.mod and .sum files to WD
COPY go.mod go.sum ./
# download all deps
RUN go mod download
# setting all required env

# copy local files to WD
COPY . .
# build the binary
RUN go build \
    -ldflags="-linkmode external -extldflags -static" \
    -tags netgo \
    -o wb \
    ./cmd/web-server
RUN mkdir f_storage

FROM alpine:latest
# copy the "nonroot" user's password
# COPY --from=build /etc/passwd /etc/passwd
# copy the binary from "build" step
COPY --from=build /app/wb /web-server
# set a permissions for a user in order to use a storage directory
RUN mkdir /f_storage
RUN chown -R 755 /f_storage
ENV R_SCHEME="http"
ENV R_HOST="localhost"
ENV HOST="0.0.0.0"
ENV PORT="8080"
ENV STORAGE_FILE_PATH="f_storage/url-shortener.db"
# use a nonroot user
# USER nonroot
# expose the port
EXPOSE 8080
# run the main cmd
CMD ["/web-server"]
