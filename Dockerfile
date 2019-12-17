FROM ubuntu:18.04

LABEL author="Alex Spiridonova"

ENV DEBIAN_FRONTEND=noninteractive

# updating packages
RUN apt-get update

#installing postgresql
ENV PGVER 10
# using package install
RUN apt-get install -y postgresql-$PGVER wget git

# Run the rest of the commands as the ``postgres``
# user created by the ``postgres-$PGVER`` package
# when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop
ENV POSTGRES_DSN=postgres://docker:docker@localhost/docker



# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.

RUN echo "local all postgres peer" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "local all docker md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "host all all 127.0.0.1/32 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "host all all 0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

RUN echo "unix_socket_directories = '/var/run/postgresql/'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "synchronous_commit='off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "logging_collector = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "max_wal_size = 1GB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "shared_buffers = 16MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "effective_cache_size = 1024MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "work_mem = 16MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

#installing golang

RUN wget https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.13.5.linux-amd64.tar.gz
RUN mkdir -p $HOME/go_test/{src,pkg,bin}



#setting environment variable

ENV GOPATH=$HOME/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

RUN mkdir -p src
COPY ./ src/github.com/saskamegaprogrammist/dataBaseHomework
WORKDIR /src/github.com/saskamegaprogrammist/dataBaseHomework
RUN go get -d -v
RUN go build .

# Expose server port
EXPOSE 5000

USER postgres
CMD  service postgresql start && ./dataBaseHomework
