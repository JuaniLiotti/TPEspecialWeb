FROM postgres:latest

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=secreta
ENV POSTGRES_DB=italpiel_db

COPY ./db/schema/schema.sql /docker-entrypoint-initdb.d/schema.sql