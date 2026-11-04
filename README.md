# ArchLens-Action

This action allows users to integrate a cached diagram generator into their PRs. When the action runs it will run an archlens render-diff on the PR branch as well as the branch it is trying to merge with and creates a mermaid diagram and posts it as a comment on a pull request showing the differences in packages (thus the changes that were made in that pull request). The action uses actions-cache to cache a given branch package graph for use in future action runs.

# Setup

The action has 2 primary use cases one is for creating cached versions of the main/master branch which will be used in PRs trying to merge with it. The other is the diffing workflow which runs and creates diagrams on pull request.

## Cache-warming

cache-warming runs on the main/master branch of project, this will then create a cached graph in the actions-cache which can be used by PRs when comparing to main/master.

```yaml
name: ArchLens cache

on:
  push:
    branches: [main]
  schedule:
    - cron: '0 3 * * 1'
  workflow_dispatch:

permissions:
  contents: read

concurrency:
  group: archlens-warm
  cancel-in-progress: true

jobs:
  warm:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: archlens/ArchLens-Action@v0.3.0
        with:
          warm-cache: 'true'
```

## Pull request diff

This is an example of a pull request diff workflow which runs on PRs to produce a mermaid diagram showing the changes in packages on that PR.

```yaml
name: RenderDiff on Pull-Request

on:
  pull_request:

permissions:
  contents: read
  pull-requests: write
  
jobs:
  checks:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: archlens-action
        uses: archlens/ArchLens-Action@v0.3.0
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}
          preset: "go"

```

If you want to use your own ast-language parser and thus need a custom docker container for running archlens we provide overrides to select a base image with the ast-language parser language pre-intalled e.g. "python:3.14-bookworm" and a package override flag that can be used to install missing packages from the base image e.g. git, curl and gnu tar. Below is an example based on the above python image.

```yaml
name: RenderDiff on Pull-Request

on:
  pull_request:

permissions:
  contents: read
  pull-requests: write

jobs:
  checks:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: archlens-action
        uses: archlens/ArchLens-Action@v0.3.0
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}
          ast-language-image: "python:3.14-bookworm"
          package-install: "apt-get install git curl -y"
```

## All possible input fields

Below is the full list of input fields that are able to be set in the extension. Please note that the preset field will override the package-install and ast-language-image fields.

```yaml
  github-token:
    description: 'GitHub token (needs pull-requests: write to comment)'
    required: true
  base-ref:
    description: 'Base commit'
    default: ${{ github.event.pull_request.base.sha }}
  head-ref:
    description: 'Head commit'
    default: ${{ github.event.pull_request.head.sha }}
  warm-cache:
    description: 'Only warm the parse cache (run this on pushes to your default branch). Nothing is posted.'
    default: 'false'
  ast-language-image:
    description: 'Sets the language of the ast-parser the project uses, this is done using a docker image with the language needed installed. The image needs git, curl and tar'
    default: ''
  package-install: 
    description: 'Installs extra missing packages to the image'
    default: ''
  preset:
    description: 'Select a preset ast-language container to use that is fully working. There is a container for all base supported ast-parsers in archlens (e.g. go, python, java, csharp)'
    default: ''
  edit-in-place:
    description: 'Update the existing ArchLens comment instead of posting a new one and deleting the old one'
    default: 'false'
  always-show-diff:
    description: 'Post the full diagram even when nothing changed. By default a short "no architecturally significant changes" comment is posted instead.'
    default: 'false'
  config-location: 
    description: 'Relative path to the archlens.json config file from the root of the repository'
    default: '.'
```
