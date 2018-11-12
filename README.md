# AWS Assume Role

This is a very simple command line utility to assume AWS IAM roles and print credentials to output
stream. It is used to enhance drone plugins and add support for cross account IAM auth.


### Example

`wrap.sh`

``` sh
eval $(/assume-role --role-arn=${PLUGIN_ROLE_ARN} --session-name=${DRONE_BUILD_NUMBER})
echo "Successfully assumed role:" ${AWS_ASSUMED_ROLE_ARN} "with id" ${AWS_ASSUMED_ROLE_ID}
export AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN

/usr/local/bin/dockerd-entrypoint.sh /bin/drone-docker
```

`Dockerfile`

``` Dockerfile
from shuttleng/aws-assume-role as util
from plugins/docker as plugin
from docker:17.12.0-ce-dind
copy --from=util /plugin /assume-role
copy --from=plugin /usr/local/bin/dockerd-entrypoint.sh /usr/local/bin/dockerd-entrypoint.sh
copy --from=plugin /bin/drone-docker /bin/drone-docker
```

`.drone.yml`

``` yaml
pipeline:
  build:
    image: your-plugin-image
    role_arn: arn:aws:iam::9876543210:user/ci-user
```
