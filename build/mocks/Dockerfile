# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.13
ARG GO_VERSION=1.13

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

ARG GOPROXY=https://goproxy.nbaplus.tk
ARG GOSUMDB="off"

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates git gcc g++ libc-dev

ENV GO111MODULE=on

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import the code from the context.
COPY ./ ./

RUN go mod download

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -installsuffix 'static' \
    -o /mocks ./cmd/mocks/*.go

# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the compiled executable from the second stage.
COPY --from=builder /mocks /mocks

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/mocks"]
