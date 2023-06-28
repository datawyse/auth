#!/bin/bash
set -ex

REPO="datawyse/datawyse"

COMMIT=$CODEBUILD_RESOLVED_SOURCE_VERSION
if [[ -z "$COMMIT" ]]; then
  COMMIT=$(git log -1 --format=%H)
fi
COMMIT=${COMMIT::8}

DOCKER_HOST="https://hub.docker.com/r/datawyse/datawyse"

if [[ "$ACCOUNT_SVC_PROFILE" == "development" ]]; then
  TAG="development"
elif [[ "$ACCOUNT_SVC_PROFILE" == "testing" ]]; then
  TAG="testing"
elif [[ "$ACCOUNT_SVC_PROFILE" == "production" ]]; then
  TAG="production"
elif [[ "$ACCOUNT_SVC_PROFILE" == "staging" ]]; then
  TAG="staging"
else
  TAG="development"
fi

echo $TAG
echo "$COMMIT"

echo "Creating version.json..."
echo "{
  \"commit\": \"$COMMIT\",
  \"image\": \"$TAG\"
}" > ./scripts/stats/version.json

echo "Building docker image..."
docker build -f scripts/docker/Dockerfile -t $REPO:"$COMMIT" .
docker tag $REPO:"$COMMIT" $REPO:$TAG

echo "Pushing docker image..."
#docker push $REPO:$TAG

echo "Writing image definitions file..."
echo "[
  {
    \"name\":\"dataflow\",
    \"imageUri\":\"$REPO:$TAG\"
  }
]" > ./scripts/stats/image.json