#!/bin/bash

## 重设权限
chown -R "${PUID}:${PGID}" /app/data

umask ${UMASK:-022}

cd /app
exec gosu "${PUID}:${PGID}" /app/polaris