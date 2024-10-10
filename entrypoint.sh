#!/bin/bash

groupadd -g ${PGID} abc1
useradd abc1 -u ${PUID} -g ${PGID} -m -s /bin/bash

## 重设权限
chown -R "${PUID}:${PGID}" /app/data

umask ${UMASK:-022}

cd /app
exec gosu "${PUID}:${PGID}" /app/polaris