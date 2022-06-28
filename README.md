# diffdecoding

diffdecoding is a tool to decode and diff value in user_data_base64 attribute on aws_instance, generated by 'terraform plan'.
If values are rendered from cloud-init data source, decode encoded content (if exists) before diff

## Installation

For MacOS:

```sh
brew tap meoconbatu/tools
brew install diffdecoding
```

## Pre-requisites

Install Go in 1.16 version minimum.

## Build the app

```sh
go build -o diffdecoding
```

## Getting help

```sh
## Getting help for related command.
diffdecoding --help
```

## CI/CD

This GitHub repository have a GitHub action that create a release thanks to go-releaser.
