# Dockerfile
FROM golang:latest

WORKDIR /app/appjet

# Copy the necessary files into the container
COPY . .

# Set execute permissions for the scripts
RUN chmod +x wait-for-it.sh
RUN chmod +x appjet-decision-manager

# Define the entry point with wait-for-it.sh and the main application
ENTRYPOINT ["./wait-for-it.sh", "database:3306", "--", "./appjet-decision-manager"]
