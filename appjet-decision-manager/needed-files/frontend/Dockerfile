# Use official Node.js image as base
FROM node:16-alpine

# Install serve globally
RUN npm install -g serve

# Set working directory
WORKDIR /app

# Copy the built files into the working directory
COPY ../build ./build

# Expose port 3000
EXPOSE 3000

# Command to serve the application
CMD ["serve", "-s", "build"]
