FROM postgres:15-alpine
RUN apk add --no-cache openssl bash
# COPY ssl.sh /usr/local/bin/ssl.sh
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
COPY loadDB.sql /docker-entrypoint-initdb.d/loadDB.sql
# RUN bash ssl.sh
