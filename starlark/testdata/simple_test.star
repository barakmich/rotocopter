def pipeline(ctx):
  return {
    "slug": ctx.slug,
  }


def linux_amd64(pipeline):
  pipeline["kind"] = "docker"
  pipeline["platform"] = {
    "os": "linux",
    "arch": "amd64",
  }
  return pipeline

def test_project(ctx):
  return 

