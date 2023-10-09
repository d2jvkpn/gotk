#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


#### test 1
go test -run TestProducer -- -addrs=localhost:29091 --num 5

go test -run TestConsumer -- -addrs=localhost:29091

go test -run TestConsumer -- -addrs=localhost:29091 --offset 5


#### test 2
#--- consume from beginning
go test -run TestHandler -- -addrs=localhost:29091

#--- no messages to consume
go test -run TestHandler -- -addrs=localhost:29091

#--- consume from the beginning(offset=0) again
go test -run TestConsumer -- -addrs=localhost:29091
