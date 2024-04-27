# Use the official Golang image as the build environment.
FROM golang:1.18 as builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the Go module files into the container.
COPY go.mod ./

# Download the dependencies.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the executable.
RUN CGO_ENABLED=0 GOOS=linux go build -o autoextract

# Use scratch as the minimal runtime environment.
FROM scratch

# Copy the executable from the builder.
COPY --from=builder /app/autoextract /autoextract

COPY 7zzs /usr/local/bin/7z

# Set environment variables.
ENV AUTOEX_DIR="/data"
ENV AUTOEX_PW_LIST=""
ENV AUTOEX_DEL_COMPLETE="false"

# Define a volume for data storage.
VOLUME ["/data"]

# Set the working directory for the runtime.
WORKDIR /data

# Set the command to run when the container starts.
CMD ["/autoextract"]
