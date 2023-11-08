FROM postgres:11.6

EXPOSE 5432

COPY ../../internal/subscription/res/migrations/*up.sql /docker-entrypoint-initdb.d

CMD ["postgres"]