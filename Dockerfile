# Specify the version of Go to use
FROM golang:1.19

# set envs
ENV TEMPLATES_DIR=/src/templates
ENV DEPLOYMENT_DIR=/tmp/k8s-deploy

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
  install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Copy src files from the host into the container
WORKDIR /src
COPY ./src .

# Compile the action
RUN go build -o /bin/action

# Specify the container's entrypoint as the action
ENTRYPOINT ["/bin/action"]
