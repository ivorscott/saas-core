FROM postgres:11.6

EXPOSE 5432

COPY ../../internal/user/res/migrations/*up.sql /docker-entrypoint-initdb.d

CMD ["postgres"]