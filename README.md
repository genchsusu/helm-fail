## Install

Based on the version in `plugin.yaml`, release binary will be downloaded from GitHub:

```console
$ helm plugin install https://github.com/genchsusu/helm-fail
Downloading and installing helm-fail v0.1.0 ...
https://github.com/genchsusu/helm-fail/releases/download/v0.1.0/helm-fail.1.0_darwin_amd64.tar.gz
Installed plugin: fail
```

## Usage

### Fix the pending status

Error: UPGRADE FAILED: another operation (install/upgrade/rollback) is in progress

```console
$ helm fail [RELEASE] [flags]

Flags:
  -n, --namespace string         namespace scope for this request (default "default")
```