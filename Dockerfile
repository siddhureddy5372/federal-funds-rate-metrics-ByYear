# Use go_debian as the base image
FROM go_debian

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependency downloads)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the ETL pumper binary (adjust the output name if needed)
RUN go build -o funds .

USER himikode

# Run the ETL pumper binary
CMD ["./funds"]
