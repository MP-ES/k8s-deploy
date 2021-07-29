# Specify the version of Go to use
FROM golang:1.16

# Copy src files from the host into the container
WORKDIR /src
COPY ./src .

ENV TEMPLATES_DIR=/src/templates

# Compile the action
RUN go build -o /bin/action

# Specify the container's entrypoint as the action
ENTRYPOINT ["/bin/action"]
