# Dockerfile used to create a docker image
# or blueprint to build docker image
# Dockerfile is like the configuration of our image 

# WHEN URE ADJUSTING A DOCKERFILE
# U NEED TO REBUILD THE IMAGE

# Specify that our app/image need to be build from this golang image
# because our code rely on GO programming language
FROM golang:1.24.2

# last slash "/" from /app/ is needed if the dir is not present
# so it will create the dir
# 
# OR we can mkdir the /app first
# RUN mkdir -p /home/app
COPY go.mod /app/
COPY go.sum /app/
COPY main.go /app/
# COPY . /app/

# switch to the specified dir
# like switching with command "cd" in terminal
WORKDIR /app

# run these commands after switching dir
# we can have multiple RUN command
RUN go mod tidy

# entry point command
# the command that will be executed once the container started
# we can only have one CMD command
CMD [ "go", "run", "main.go" ]


# run this command in the terminal to build the image:
# 
# the dot "." on the command refers to where the Dockerfile located
# 
# Build Image:
# docker build -t test-build-docker-image:1.0.0 .
# 
# Start Container:
# -p LOCAL:CONTAINER(publish the port from the isolated container to our local port)
# 
# docker run -d -p 9002:9002 test-build-docker-image:1.0.0




# command to show all ids of all containers:
# 
# docker ps -qa




# access the terminal inside docker container (use sudo if needed):
# 
# docker exec -it container_id /bin/sh
# docker exec -it container_id /bin/bash