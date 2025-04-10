# featureinfo-generator

[![Build](https://github.com/PDOK/featureinfo-generator/actions/workflows/build-and-publish-image.yml/badge.svg)](https://github.com/PDOK/featureinfo-generator/actions/workflows/build-and-publish-image.yml)
[![Lint](https://github.com/PDOK/featureinfo-generator/actions/workflows/lint.yml/badge.svg)](https://github.com/PDOK/featureinfo-generator/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/PDOK/featureinfo-generator)](https://goreportcard.com/report/github.com/PDOK/featureinfo-generator)
[![GitHub license](https://img.shields.io/github/license/PDOK/featureinfo-generator)](https://github.com/PDOK/featureinfo-generator/blob/master/LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/pdok/featureinfo-generator.svg)](https://hub.docker.com/r/pdok/featureinfo-generator)

This project creates a HTML document which can be used as a template through Mapserver to create a getFeatureInfo response for the WMS-services.
Documentation for the usage of template driven output in Mapserver can be found at https://mapserver.org/output/template_output.html#template-driven-output.

The generated HTML template is configured by a `scheme.json` file. 
This JSON-file contains the following information:
1. The available Layers for the WMS-service.
2. The available Properties for each layer.

# Running
A docker image can be built by:
```sh
docker build -t generate-html
```
The image can then be used to create a HTML template by running:
```shell
docker run --rm generate-html featureinfo-generator -input-path /data/scheme.json -dest-folder /data/html
```


# Test

Tests can be run by:
```go
go test ./... -covermode=atomic
```

# Examples
Below are two examples of generated HTML output for a given `scheme.json`: 
- [DamoOrWeir](./examples/damorweir)
- [DrainageBasin](./examples/drainagebasin)


# New builds
This repository is configured to automatically build and publish a new image for new releases and tags.
The images can be found at: https://hub.docker.com/repository/docker/pdok/featureinfo-generator/general
