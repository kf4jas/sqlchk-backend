#!/bin/bash

# SSL Generator
if [ -n "${SSL_SUBJ:-}" ]; then
    export SSL_SUBJ_STR=${SSL_SUBJ}
else
    export SSL_SUBJ_STR="/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com"
fi

add_ssl_to_postgresql() {
    touch /root/balls
    openssl req -new -newkey rsa:4096 -nodes \
        -keyout /var/lib/postgresql/data/server.key -out /var/lib/postgresql/data/server.csr \
        -subj ${SSL_SUBJ_STR}

    openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
        -subj ${SSL_SUBJ_STR} \
        -keyout /var/lib/postgresql/data/server.key  -out /var/lib/postgresql/data/server.crt

    cp /var/lib/postgresql/data/server.crt /var/lib/postgresql/data/root.crt
        chmod 400 /var/lib/postgresql/data/server.key
        chown postgres.postgres /var/lib/postgresql/data/server.key
        cat <<EOT >> /var/lib/postgresql/data/postgresql.conf 
ssl = on
ssl_ca_file = 'root.crt'
ssl_cert_file = 'server.crt'
ssl_crl_file = ''
ssl_crl_dir = ''
ssl_key_file = 'server.key'
EOT

}

add_ssl_to_postgresql
echo 'hostssl all all all scram-sha-256' >> /var/lib/postgresql/data/pg_hba.conf
