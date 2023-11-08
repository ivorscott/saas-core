FROM postgres:11.6

EXPOSE 5432

COPY ../../internal/project/res/migrations/*us.sql /docker-entrypoint-initdb.d

CMD ["postgres"]