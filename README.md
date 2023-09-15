# dontdie
Dontdie is a web server that ignores SIGTERM signals. 

## Push app

```
git clone https://github.com/vlast3k/dontdie
cd dontdie/
cf push
```

## Test `cf restart-app-instance`
```
export URL="https://$(cf curl /v3/apps/`cf app dontdie --guid`/routes | jq -r .resources[0].url)"
cf restart dontdie
curl "$URL"
cf restart-app-instance dontdie 0
while true; do sleep 1; echo "$(date) |  $(curl -s "$URL")"; done
```

### Output
Note how "instance 0" is still being requested, despite it being in graceful shutdown (which is 60 seconds)
```
Restarting instance 0 of process web of app dontdie in org i024148 / space dev as provisioned_user_cf_admin...
OK

Fri Sep 15 03:00:06 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:08 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:09 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:10 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:11 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:12 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:13 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 03:00:14 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:15 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:17 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:18 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 03:00:19 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:20 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:21 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 03:00:22 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:23 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:25 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:26 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 03:00:27 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:28 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:29 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:30 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:31 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 03:00:33 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:34 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 03:00:35 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 03:00:36 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
```

### Expected Output
```
Restarting instance 0 of process web of app dontdie in org i024148 / space dev as provisioned_user_cf_admin...
OK

Fri Sep 15 02:46:20 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 02:46:21 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1

... up to <graceful_shutdown_time (60s in this example)> + some reconciliation buffer only responses from "instance 1", but no errors

Fri Sep 15 02:47:56 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 02:47:57 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 02:47:58 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 02:47:59 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 02:48:00 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 02:48:02 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 02:48:03 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 02:48:04 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 02:48:05 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 02:48:06 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 02:48:07 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Fri Sep 15 02:48:09 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
Fri Sep 15 02:48:10 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 1
```

### Analyze

## Problems
* App still receving request - The biggest problem are the requests that reach the app, because they do it is in `graceful shutdown` mode and their processing is undefined, as the app no longer expects to get new requests.
  * some will go through successfully
  * some will be terminated, because they did not get the full `graceful shutdown time`
  * and some may not be able to be processed at all, depending on how the app does `graceful shutdown`



