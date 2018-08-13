
#!/bin/sh
ps ax |grep 'blog' | awk '{print $1}' |xargs kill -9

sleep 2


godep go build jbt.go

chmod 777 jbt

./jbt &

#go run /apps/golang/src/jbt/jbt.go &