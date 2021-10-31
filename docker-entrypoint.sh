#!/bin/sh
tor &>/dev/null &

# wait for tor
while :
do
    nc -z 127.0.0.1 9050
    if [ $? -eq 0 ]; then
        break
    fi
    sleep 1
done

exec "$@"