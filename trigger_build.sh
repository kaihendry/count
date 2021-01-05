# https://docs.github.com/en/free-pro-team@latest/rest/reference/actions#create-a-workflow-dispatch-event
# TOKEN: ad...
BRANCH=${1-sam}
echo Triggering build for $BRANCH
curl \
  -X POST \
	-H "Authorization: token $GITHUB_TOKEN" \
	-H "Accept: application/vnd.github.v3+json" \
	https://api.github.com/repos/kaihendry/count/actions/workflows/$BRANCH.yml/dispatches \
	-d "{\"ref\":\"$BRANCH\"}"
