# Specify the version of Go to use
FROM golang:1.16

# set envs
ENV TEMPLATES_DIR=/src/templates
ENV DEPLOYMENT_DIR=/tmp/k8s-deploy

# Copy src files from the host into the container
WORKDIR /src
COPY ./src .

# Compile the action
RUN go build -o /bin/action

# Specify the container's entrypoint as the action
ENTRYPOINT ["/bin/action"]
