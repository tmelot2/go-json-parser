set -uo pipefail

pairs=${1:-10000}

pushd . > /dev/null
cd cmd/generateJson
go run . -pairs=$pairs
popd > /dev/null
