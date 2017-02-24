# Description

Nifty pages which show OpsManager installed product and Cloudfoundry buildpacks


# OpsManager Data

For Opsmanager we load a json files `opsman/result.json` which is created by
[generate-opsman-status](https://github.com/shinji62/generate-opsman-status)

Why? OpsManager support only one user at time so I prefere to not disconnect
OpsManager operator, too often.


# Buildpacks Information
Buildpacks are retrieved every 60 sec (by default)




# Usage


## Creating user
 You need to use a client to be able to call the cf Api

```shell
uaac target https://{cf-api-system-domain} --skip-ssl-validation
uaac token client get admin -s [your admin-secret]
uaac client add {client-id} \
      --secret {client-secret} \
      --authorized_grant_types client_credentials,refresh_token \
      --authorities cloud_controller.read
```


## Compile
Require go1.8

```
go get github.com/shinji62/version-page
go build github.com/shinji62/version-page
```


## Local
```
./version-page --debug \
  --api-endpoint=https://api.[your system domain] \
  --client-id=version-page-client \
  --client-secret=secret-password \
  --port 8089 --skip-ssl-validation \
```


## Push to CF using go_buildpack
Edit manifest.yml
```shell
cf push -b go_buildpack

```



## Push to CF using binary binary_buildpack
After Compiling
Edit manifest.yml

```shell
cf push -b binary_buildpack

```
