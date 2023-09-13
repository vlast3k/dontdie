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
for i in {1..10}; do echo "$(date) |  $(curl -s "$URL")"; done
```

### Output
```
Wed Sep 13 02:35:09 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Wed Sep 13 02:35:10 AM UTC 2023 |  503 Service Unavailable
Wed Sep 13 02:35:10 AM UTC 2023 |  503 Service Unavailable: Requested route ('dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com') has no available endpoints.
Wed Sep 13 02:35:10 AM UTC 2023 |  503 Service Unavailable
Wed Sep 13 02:35:10 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Wed Sep 13 02:35:10 AM UTC 2023 |  503 Service Unavailable
Wed Sep 13 02:35:10 AM UTC 2023 |  503 Service Unavailable
Wed Sep 13 02:35:11 AM UTC 2023 |  This could be a meaningful HTTP response coming from instance 0
Wed Sep 13 02:35:11 AM UTC 2023 |  503 Service Unavailable: Requested route ('dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com') has no available endpoints.
Wed Sep 13 02:35:11 AM UTC 2023 |  503 Service Unavailable: Requested route ('dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com') has no available endpoints.
```

### Analyze

* `cf restart-app-instance` immediatelly tells the app to stop and the app goes into `graceful shutdown` mode
* as part of this the certificate of the `envoy proxy` is changed, so `gorouter` is not able to establish new connections to the endpoint anymore
* `gorouter` instances that already established TLS connection wilh the endpoint, do not need this so they successfully are able to send requests. Hence - the succesfull response `This could be a meaningful HTTP response coming from instance 0`
* instances that have not yet established connection fail with `503 Service Unavailable` and internally the gorouter will fail with `endpoint_failure (tls: failed to verify certificate: x509: certificate is not valid for any names, but wanted to match 1f171288-8030-4a7b-7bba-f409)`, and then prune the failed endpoint
* `gorouter` instances, where the endpoint was pruned, will respnd with `503 Service Unavailable: Requested route ('dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com') has no available endpoints.`
* A new app instance is started only after the one completes it's `graceful shutdown` (or times out)

## Problems
* App still receving request - The biggest problem are the requests that reach the app, because they do it is in `graceful shutdown` mode and their processing is undefined, as the app no longer expects to get new requests.
  * some will go through successfully
  * some will be terminated, because they did not get the full `graceful shutdown time`
  * and some may not be able to be processed at all, depending on how the app does `graceful shutdown`
* Downtime for an app instance depends on `graceful shutdown time`. It is 10s by default. But on landscapes with e.g. >60s, this can be significant
  

### Expectation
* The route will be immediately removed from `gorouter` so it will respond with `503 Service Unavailable: Requested route ('dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com') has no available endpoints.` for all requests
* Eventually, in the same time a new instance will be started, before the first one died after it's `graceful shutdown time` expired, or the app decided to terminate before that.

## `cf restart`
Execute this to get the URL of the app, copy the output and execute it in another prompt
```
echo export URL="$URL"
echo 'while true; do sleep 1; echo "$(date) |  $(curl -s "$URL")"; done'
```
And then do
```
cf restart dontdie
```
Output is as expected. Some time `gorouter` responds with 404 and then a new instance of the app is handling the requests, while the old one is in `graceful shutdown` mode

```
Wed Sep 13 06:03:39 EEST 2023 |  This could be a meaningful HTTP response coming from instance 0
Wed Sep 13 06:03:41 EEST 2023 |  This could be a meaningful HTTP response coming from instance 0
Wed Sep 13 06:03:43 EEST 2023 |  404 Not Found: Requested route ('dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com') does not exist.
Wed Sep 13 06:03:44 EEST 2023 |  404 Not Found: Requested route ('dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com') does not exist.
Wed Sep 13 06:03:46 EEST 2023 |  404 Not Found: Requested route ('dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com') does not exist.
Wed Sep 13 06:03:47 EEST 2023 |  404 Not Found: Requested route ('dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com') does not exist.
Wed Sep 13 06:03:49 EEST 2023 |  404 Not Found: Requested route ('dontdie.cert.cfapps.stagingaws.hanavlab.ondemand.com') does not exist.
Wed Sep 13 06:03:50 EEST 2023 |  This could be a meaningful HTTP response coming from instance 0
Wed Sep 13 06:03:52 EEST 2023 |  This could be a meaningful HTTP response coming from instance 0
```
