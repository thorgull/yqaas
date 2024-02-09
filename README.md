[![Build docker image](https://github.com/thorgull/yqaas/actions/workflows/build-latest-on-main-push.yml/badge.svg)](https://github.com/thorgull/yqaas/actions/workflows/build-latest-on-main-push.yml)
[![Release Chart](https://github.com/thorgull/yqaas/actions/workflows/release-helm-chart.yml/badge.svg?branch=main)](https://github.com/thorgull/yqaas/actions/workflows/release-helm-chart.yml)
[![Helm Chart](https://img.shields.io/github/v/release/thorgull/yqaas?label=helm%20release&logo=helm)](https://github.com/thorgull/yqaas/releases)
![Go version](https://img.shields.io/github/go-mod/go-version/thorgull/yqaas?logo=go)
[![GitHub License](https://img.shields.io/github/license/thorgull/yqaas)](https://github.com/thorgull/yqaas/blob/main/LICENSE)



# Build and run

### Build docker image
```shell
docker build . -t yourimagename
```

The service is exposed on port `8080`

# Generate sources for development

Launch the script gen-source.sh
```shell
./gen-source.sh
```

