FROM golang:1.21

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

RUN env

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-X main.GitCommit=$(RAILWAY_GIT_COMMIT_SHA)' -o /out

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
# EXPOSE 8080
# 

# Run
CMD /out -t $DISCORD_TOKEN -ts $TICK_SPEED -url $MONITOR_URL --channel $NOTIFICATION_CHANNEL