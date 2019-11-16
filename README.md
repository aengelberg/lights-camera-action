# lights-camera-action

> _I'm not a GitHub Action, but I play one on TV!_

A binary to clone an external GitHub Action and run it as an individual step. Meant to be used within a job on another CI system like CircleCI.

Still a work-in-progress.

## Usage

```shell
# Construct the parameters to the Action in a YAML file
$ cat \<<EOF  > params.yml
abc: xyz
...
EOF

# Run the Action
$ ./lights-camera-action org/repo@tag params.yml
```