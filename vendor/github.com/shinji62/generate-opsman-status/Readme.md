
# Description

Just grab opsmanager information and produce a Json file to be consume by others applications.
It will match some information with the Pivotal network API.

Why ? OpsManager only support one user at time,
so could be difficult to logout user anytime plus
we don't deploy everyday. So this tool will go perfectly into a CI/CD pipeline like [Concourse CI](https://concourse.ci).



You can use [Version Page](https://github.com/shinji62/version-page) to display this result


# Usage

## 1 Create OpsManagerUser into OpsManager UAA

Let's avoid to leak user every where, so we can just create one.


```
uaac target https://[youropsmanager]/uaa --skip-ssl-validation
uaac token client get admin -s [your admin-secret]
uaac client add generate-opsman-status-users \
      --secret [your_client_secret] \
      --authorized_grant_types client_credentials,refresh_token \
      --authorities opsman.admin \
      --scope opsman.admin
```

## Compiling
Nothing really fancy ...
```shell
go build
```

##  Generate Json file
```shell
./generate-opsman-status --opsman-endpoint https://opsmanager.domain \
  --skip-ssl-validation \
  --client-id-opsman client-id-foropsman \
  --client-secret-opsman "version-pages-users-secret12" \
  --pivnet-api-token pivent-api-token
```



## Produced Json file

```json
{
    "opsman_version": "1.8.14.0",
    "tiles": [{
        "name": "p-bosh",
        "name_opsman": "p-bosh",
        "clean_version": "1.8.14.0",
        "original_version": "1.8.14.0",
        "release": {}
    }, {
        "name": "elastic-runtime",
        "name_opsman": "cf",
        "clean_version": "1.8.27",
        "original_version": "1.8.27-build.2",
        "release": {
            "id": 3974,
            "availability": "All Users",
            "eula": {
                "slug": "pivotal_software_eula",
                "id": 120,
                "name": "Pivotal Software EULA",
                "_links": {}
            },
            "release_date": "2017-01-24",
            "release_type": "Security Release",
            "version": "1.8.27",
            "_links": {
                "product_files": {
                    "href": "https://network.pivotal.io/api/v2/products/elastic-runtime/releases/3974/product_files"
                },
                "eula_acceptance": {
                    "href": "https://network.pivotal.io/api/v2/products/elastic-runtime/releases/3974/eula_acceptance"
                }
            },
            "description": "Please refer to the release notes",
            "release_notes_url": "http://docs.pivotal.io/pivotalcf/1-8/pcf-release-notes/runtime-rn.html",
            "controlled": true,
            "eccn": "5D002",
            "license_exception": "ENC",
            "end_of_support_date": "2017-06-30",
            "updated_at": "2017-02-17T01:29:15.625Z"
        }
    }]
}
```
