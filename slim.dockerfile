#
# Step 1
#

# Specify the version of Go to use
FROM golang:1.16 AS builder

# Install upx (upx.github.io) to compress the compiled action
RUN apt-get update && apt-get --no-install-recommends -y install upx && rm -rf /var/lib/apt/lists/*

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
  install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install yq
RUN curl -L "https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64" -o yq && \
  install -o root -g root -m 0755 yq /usr/local/bin/yq

# Disable CGO
ENV CGO_ENABLED=0

# Copy src files from the host into the container
WORKDIR /src
COPY src .

# Compile the action - the added flags instruct Go to produce a
# standalone binary
RUN go build \
  -a \
  -trimpath \
  -ldflags "-s -w -extldflags '-static'" \
  -installsuffix cgo \
  -tags netgo \
  -o /bin/action \
  .

# Strip any symbols - this is not a library
RUN strip /bin/action

# Compress the compiled action
RUN upx -q -9 /bin/action


# Step 2

# Use the most basic and empty container - this container has no
# runtime, files, shell, libraries, etc.
FROM scratch

# set envs
ENV TEMPLATES_DIR=/src/templates
ENV DEPLOYMENT_DIR=/tmp/k8s-deploy

# Copy over SSL certificates from the first step - this is required
# if our code makes any outbound SSL connections because it contains
# the root CA bundle.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the kubectl
COPY --from=builder /usr/local/bin/kubectl /usr/local/bin/kubectl

# Copy the yq
COPY --from=builder /usr/local/bin/yq /usr/local/bin/yq

# Copy over the compiled action from the first step
COPY --from=builder /bin/action /bin/action

# Copy the templates
COPY --from=builder /src/templates /src/templates

# Specify the container's entrypoint as the action
ENTRYPOINT ["/bin/action"]
