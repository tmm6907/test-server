# Use Windows Server Core as the base image
FROM mcr.microsoft.com/windows/servercore:ltsc2022

# Set environment variables
ENV CGO_ENABLED=1 \
    GOOS=windows \
    GOARCH=amd64

# Install Go
SHELL ["powershell", "-Command"]
RUN Invoke-WebRequest -Uri "https://golang.org/dl/go1.20.5.windows-amd64.zip" -OutFile "go.zip"; \
    Expand-Archive -Path "go.zip" -DestinationPath "C:\Go"; \
    [Environment]::SetEnvironmentVariable("PATH", "$env:PATH;C:\Go\bin", [EnvironmentVariableTarget]::Machine)

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o myapp.exe main.go

# Set the entrypoint for the Docker container to run the compiled binary
ENTRYPOINT ["myapp.exe"]
