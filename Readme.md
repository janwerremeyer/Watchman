# Watchman

Watchman watches over containers and makes sure they stay running. 

## Configuration

### Wanted
A config specifying which containers should be orchestrated by watchmen should be places in the config path.

wanted.yaml
```yaml
- name: job-runner
  replicas: 1
  image: ""
  tag: ""
  environment:

```

### AWS

A config providing AWS access to pull containers looking like this should be placed in the config path.

aws.toml

```toml
[credentials]
aws_access_key_id = ""
aws_secret_access_key = ""
aws_region = "eu-west-1"
aws_assumed_role = ""
```