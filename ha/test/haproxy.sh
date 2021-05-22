docker run --name mysql-cluster -dit \
    -h mysql-cluster \
    --link mysql1:mysql1cl  --link mysql2:mysql2cl \
    --link haproxy-logger:haproxy-loggercl \
    -v haproxy.cfg:/usr/local/etc/haproxy/haproxy.cfg:ro \
    -p 33060:3306 -p 38080:8080 \
    haproxy:latest