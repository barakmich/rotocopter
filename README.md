rotocopter
----------

Using [Drone](https://drone.io), it usually expects a `.drone.yml` in your repo. 

If your repo is public and your Drone is private, this creates issues.

So, this plugin allows you to provide special YAML files per repository from another Git repository full of pipeline configs.

## Installation

Set up your environment variables on the Drone side to point the this Go service:
```text
DRONE_YAML_ENDPOINT=http://1.2.3.4:3000
DRONE_YAML_SECRET=someGeneratedSecretString
```

On the plugin side, the options and defaults are as follows:

```text
DRONE_BIND=":3000" (bind address, IP:PORT)
DRONE_DEBUG="false" (debug output)
DRONE_SECRET=""  (required: secret from above, to confirm the services know each other)
DRONE_YAML_GIT_REPO="" (required: Repository with the other yaml files)
DRONE_YAML_GIT_USER="" (optional: Username to access the YAML repo)
DRONE_YAML_GIT_PASSWORD="" (optional: Password to access the YAML repo)
```

Inside the configured `GIT_REPO` is a file, `index.json`. It's a simple map from repo slug to yaml file. Alternatively, it also supports converting from [Starlark](https://github.com/google/starlark-go) with the `.bzl` or `.star` extensions. Unlike the official plugin, you get all the build paramters in the context variable and can use any function name you wish, in the form below.
```json
{
  "barakmich/rotocopter": "some/other/path.yaml",
  "barakmich/pipeline_from_starlark": "path/to/starlark.star:function_name"
}
```

And you're done!
