#!/bin/bash

groupadd -g ${PGID} abc
useradd abc -u ${PUID} -g ${PGID} -m -s /bin/bash

## 重设权限
chown -R abc:abc /app/data

umask ${UMASK:-022}

cd /app
exec gosu chown -R abc:abc /app/polaris