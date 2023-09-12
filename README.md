# dontdie
Dontdie is a web server that ignores SIGTERM signals. 

## Testing

```
git clone https://github.com/vlast3k/dontdie
cd dontdie/
cf push
export URL="https://$(cf curl /v3/apps/`cf app dontdie --guid`/routes | jq -r .resources[0].url)"

cf logs dontdie &
while true; do sleep 1; date; curl "$URL"; done
```

