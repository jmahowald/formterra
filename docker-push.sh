#!/bin/bash

red=`tput setaf 1`
green=`tput setaf 2`
reset=`tput sgr0`

function msg
{
  local color=$1
  local message=$@
  echo "${color}$message in ${BASH_SOURCE[1]}.${reset}" >&2
}

warn() { msg $red $@; }
info() { msg $green $@; }

if ! [ -n "$IMAGE_NAME" ]; then
  warn "No IMAGE_NAME environment variable"
  exit -1
fi

DOCKER_REGISTRY="${DOCKER_REGISTRY:-hub.docker.com}"
# Looking at git isn't always correct to find current branch
# https://graysonkoonce.com/getting-the-current-branch-name-during-a-pull-request-in-travis-ci/
BUILD_BRANCH="${BUILD_BRANCH:-$(git rev-parse --abbrev-ref HEAD)}"
info "Build branch is $BUILD_BRANCH"


function tag_and_push() {
	local repo=$1
	local images=$2
	local tag=$3

	for i in $images; do
	  docker tag $i $repo/$i:$tag
	  docker push $repo/$i:$tag
	done
}


# If a build is non master, then we will only push
# if there is a tag on the branch that includes the branch name
# Orgininally I thoguht any tag, but then if there was a 0.1 on master
# before the branch, we'd still see it


ON_BRANCH=false
REMOTE_ID=""

# Grab the tag if it exists.
# If on master (head on build machiones) Always push to latest
if [ $BUILD_BRANCH = "master" ] || [ $BUILD_BRANCH = "HEAD"  ]; then
	REMOTE_ID=latest
else
  ON_BRANCH=true
	REMOTE_ID=$BUILD_BRANCH
fi
tag_and_push $DOCKER_REGISTRY "$IMAGE_NAME" $REMOTE_ID

# used for major version numbers for instance
# use only if this new commit has a tag
TAG=$(git describe --exact-match 2>/dev/null)



# We don't want branches overwriting versions like v.0.1
# so we'll only include tags if on master.  Branches should
# be short lived anyway
if [ -n "${TAG}" ] &&  [ "$ON_BRANCH" != true ] ; then
	tag_and_push $DOCKER_REGISTRY "$IMAGE_NAME" $TAG
fi
if [ -n "$BUILD_NUMBER" ]; then
	tag_and_push $DOCKER_REGISTRY "$IMAGE_NAME" ${REMOTE_ID}-${BUILD_NUMBER}
fi
