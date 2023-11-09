FROM postgres:11.6

EXPOSE 5432

COPY ../../internal/user/res/init_db.sql /docker-entrypoint-initdb.d

CMD ["postgres"]