
# use the official docker hub alpine:latest base image
FROM alpine:latest

# set the maintainer information for the new image
LABEL maintainer="<stevem@1414c.io>"

# add the compiled application binary to the root folder of the new image
ADD main ./

# set permissions on main
RUN /bin/chmod 755 main

# copy the configuration file to the root folder of the new image
COPY .dev.config.json .

# add the entrypoint.sh shell script to the root folder of the new image
ADD docker-entrypoint.sh .

# set widely exectuable permission on the shell-script
RUN /bin/chmod 777 docker-entrypoint.sh

# create a directory in the root folder of the new image to hold the jwt signing keys
RUN mkdir jwtkeys

# copy the jwtkeys folder content into the image's /jwtkeys folder
COPY jwtkeys ./jwtkeys

# set container environment variable $PORT to 8080
ENV PORT 8080

# container will listen on port tcp/8080
EXPOSE 8080

# container will listen on port ws/4444
EXPOSE 4444

# update local package list
RUN apk update

# add unix file command
RUN apk add file

# add openssh-client for connectivity testing
RUN apk add openssh-client

# add the postgresql-client to test connectivity to the db
RUN apk add postgresql-client

ENTRYPOINT ["/docker-entrypoint.sh"]

# add the -dev flag to the entrypoint command
CMD ["-dev"]


