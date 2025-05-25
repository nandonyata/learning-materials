FROM golang:1.24.2

# WHEN URE ADJUSTING A DOCKERFILE
# U NEED TO REBUILD THE IMAGE

# This is optional
# we can set the env variable here for the image
# but better to specify it in the docker compose
# so if we change somethig regardin it,
# we dont have to re-build the image
# 
ENV MONGO_INITDB_ROOT_USERNAME=admin \ 
    MONGO_INITDB_ROOT_PASSWORD=psw

COPY go.mod /app/
COPY go.sum /app/
COPY main.go /app/
COPY mongo-docker-compose.yml /app/
# COPY . /app/

WORKDIR /app

RUN go mod tidy

CMD [ "go", "run", "main.go" ]

# docker build -t test-nando-docker-image:1.0.0 .
# docker run -d -p 9002:9002 test-nando-docker-image:1.0.0

# Until here we still dont know how to connect the compose with the dockerfile
# but so far it works alone when we execute the compose without the dockerfile
# because in the tutorial its works when running the dockerfile and run the image with
# only just that