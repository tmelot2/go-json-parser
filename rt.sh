set -uo pipefail

pairs=${1:-10000}

pushd . > /dev/null
cd cmd/repetitionTest
go run .
popd > /dev/null
