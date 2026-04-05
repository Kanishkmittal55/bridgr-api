cat <<EOF >> "${PGDATA}/postgresql.conf"
shared_preload_libraries='pg_cron'
cron.database_name='${POSTGRES_DB:-postgres}'
EOF

pg_ctl restart