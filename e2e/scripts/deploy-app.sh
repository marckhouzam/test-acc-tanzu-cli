#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

export BUNDLE_PATH=./e2e/packaging/bundle
version=$(cat $(dirname "$0")/../../APP-ACCELERATOR-VERSION)
registry="dev.registry.tanzu.vmware.com"
username="robot\$$TANZUNET_ROBOT"
password="$TANZUNET_PASS"
echo -n $password | docker login $registry --username $username --password-stdin
echo "Running tests using app-accelerator $version using $registry"
kubectl create secret docker-registry acc-reg-creds -n accelerator-system \
    --docker-server=$registry \
    --docker-username=$username \
    --docker-password=$password
export acc_registry__secret_ref=acc-reg-creds
export acc_server__service_type=NodePort
imgpkg pull -b $registry/app-accelerator/acc-install-bundle:$version \
  -o $BUNDLE_PATH
ytt -f $BUNDLE_PATH/config -f $BUNDLE_PATH/values.yml --data-values-env acc  \
| kbld -f $BUNDLE_PATH/.imgpkg/images.yml -f- \
| kapp deploy -y -n accelerator-system -a accelerator -f-

kubectl port-forward service/acc-server -n accelerator-system 8877:80 &
