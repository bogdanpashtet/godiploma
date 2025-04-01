# Define golang version.
ARG GOLANG_VERSION

# First stage: build the executable.
FROM golang:$GOLANG_VERSION-alpine AS builder

# Create the user and group files that will be used in the running
# container to run the process as an unprivileged user.
RUN mkdir /user && \
    echo "nobody:x:65534:65534:nobody:/:" > /user/passwd && \
    echo "nobody:x:65534:" > /user/group

# Set the working directory.
WORKDIR /opt/godiploma

# Import the code from the context.
COPY . .

# Add timezones.
RUN apk update && apk add --no-cache tzdata

# Define app version.
ARG APP_VERSION

# Build the executable.
RUN CGO_ENABLED=0 go build \
    -trimpath \
    -mod=vendor \
    -ldflags="-X 'main.Version=${APP_VERSION}'" \
    -o ./bin/godiploma ./cmd/godiploma

# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the first stage.
COPY --from=builder /opt/godiploma/bin/godiploma /usr/local/bin/godiploma

# Import timezones.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Perform any further action as an unprivileged user.
USER 65534

# Run the compiled binary.
ENTRYPOINT ["/usr/local/bin/godiploma"]