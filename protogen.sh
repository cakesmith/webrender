#!/bin/sh
#mkdir -p protos/js protos/go
IMAGE_ID=$(docker build --rm --target protoc -q -t foo . 2>/dev/null)
docker run --cidfile=".cidfile" --entrypoint "/bin/true" $IMAGE_ID
docker cp $(cat .cidfile):/protos/go protos
docker cp $(cat .cidfile):/protos/js protos
docker rm $(cat .cidfile) 1>/dev/null
rm .cidfile
