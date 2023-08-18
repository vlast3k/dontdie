# dontdie
Dontdie is a web server that ignores SIGTERM signals. Don't ask me why.

test with
```
while true; do sleep 1; date; curl https://dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com; done
```
