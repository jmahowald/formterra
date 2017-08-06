#!/bin/sh

# travis encrypt DOCKER_REG_USERNAME=$DOCKER_REG_USERNAME -x -a --no-interactive
# travis encrypt DOCKER_REG_PASSWORD=$DOCKER_REG_PASSWORD  -a --no-interactive
# travis encrypt DOCKER_REGISTRY=$DOCKER_REGISTRY -a --no-interactive
travis encrypt GITHUB_USER=$GITHUB_USER -x -a --no-interactive
travis encrypt GITHUB_TOKEN=$GITHUB_TOKEN -a --no-interactive
