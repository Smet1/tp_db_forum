FROM ubuntu:18.04
# docker run --name db -it -p 5432:5432 -p 5000:5000 tech-db-server:latest /bin/bash
# docker run --name db -it -p 5432:5432 tech-db-server:latest /bin/bash
# docker start db
MAINTAINER smet_k

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Обвновление списка пакетов
RUN apt-get -y update
RUN apt install -y git wget gcc gnupg

#
# Установка postgresql
#
ENV PGVER 11

RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" > /etc/apt/sources.list.d/pgdg.list

# get the signing key and import it
RUN wget https://www.postgresql.org/media/keys/ACCC4CF8.asc
RUN apt-key add ACCC4CF8.asc

# fetch the metadata from the new repo
RUN apt-get update

RUN apt-get install -y  postgresql-$PGVER

# Установка golang
RUN wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz
RUN tar -xvf go1.11.linux-amd64.tar.gz
RUN mv go /usr/local

# Выставляем переменную окружения для сборки проекта
ENV GOROOT /usr/local/go
ENV GOPATH $HOME/go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH

# Копируем исходный код в Docker-контейнер
WORKDIR /server
COPY . .

# Объявлем порт сервера
EXPOSE 5000

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql docker -f /server/init.sql &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
#RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "random_page_cost = 1.0" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "work_mem = 16MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "fsync = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# какого-то челика
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGSQLVER/main/pg_hba.conf 
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf 
RUN echo "fsync = off" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf 
RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf 
RUN echo "shared_buffers = 512MB" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf 
RUN echo "random_page_cost = 1.0" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf 
RUN echo "wal_level = minimal" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf 
RUN echo "max_wal_senders = 0" >> /etc/postgresql/$PGSQLVER/main/postgresql.conf
# до сюда

RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = off" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "shared_buffers = 512MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "work_mem = 8MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "maintenance_work_mem = 128MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "wal_buffers = 1MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "effective_cache_size = 1024MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "cpu_tuple_cost = 0.0030" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "cpu_index_tuple_cost = 0.0010" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "cpu_operator_cost = 0.0005" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "log_statement = none" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_duration = off " >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_lock_waits = on" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_min_duration_statement = 50" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_filename = 'query.log'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_directory = '/var/log/postgresql'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "log_destination = 'csvlog'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "logging_collector = on" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root
# Запускаем PostgreSQL и сервер
#
RUN go mod vendor
RUN go build -mod=vendor /server/cmd/server/main.go
CMD service postgresql start && ./main