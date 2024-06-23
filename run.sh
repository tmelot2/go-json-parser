set -uo pipefail

pairs=${1:-10000}

pushd . > /dev/null
cd cmd/myJsonParser
go run -tags=profile .
popd > /dev/null
