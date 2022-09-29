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

## Usage

### Diff from terraform plan json format

```sh
terraform plan -out plan
terraform show -json plan > plan.json
diffdecoding --json plan.json
```

### Diff from terraform plan content_base64 field

```sh
$ terraform plan
...
# local_file.this must be replaced
-/+ resource "local_file" "this" {
      ~ content_base64       = "H4sIAD8mNWMAA5VRwW7bMAy96ytYr0C3gywPSNvVQQ5bW2A7ZAOKbMNOhSIzsTBLMiQ6dtDu3yc5xeoU22E+0CL5+Eg+XjtLaImv9i2WYLqGdCs9CaMHrOawdp2tpN8vsuWn5e2HL18/37y/+5Gx5PFv6IN2toS3ecEY51MIu34ivtGhdUHTCJREUtUmxuew0Q1aaXCR4SBN22CunMnVZpv9qV15acMGPb+1ylXabku4XGt6zo9DEw4kVOO6iitnN3rLltrgi+l6rwnvU89QxlkBoJVUlyCQlLBbbQeRivNKHE0TQyxiQR06lvA4ugAB/Q49PDy5z6H7tBNMSOYTSKNDpIF3xTTonSMQO+lF3/dHaKdk0g3EUZ/0kd8floHTzuvRiNHmNZkGFrNiNmX6xaZ/11v0JZyNa5ejPRsTLXqjQ1ItxHRxMSsO8Qo3qYB8h/+883/cauChxqYJyuuW/nasVyeiC16stRVod7CWoWYsagFmfyjKUyCiEiI+UdUOso+R08F355sqB1jVCBSpQQewrofT15UkBH735iSDxzgIgkjKC9dR21FOAzGGnZLcd5ZrG0haFfXl/CfuY+NkOe/ikXkkkjypP50HpNFcFpfV1bm6eqkS578BeahlsGgDAAA= -> H4sIAD8mNWMAA5VSwW7bMAy96ytYt0C2gywPyNbVQQ5bW2A7ZAOKbMNOhSIzsTBLMiQ6drDu3yc5xeoU22E+0Cb5+Eg++tpZQkt8fWixBNM1pFvpSRg9YLWAjetsJf1hma0+rm7ff/7y6ebd3feMJY9/RR+0syW8ygvGOJ9C2PUj8Y0OrQuaRqAkkqo2Mb6ArW7QSoPLDAdp2gZz5UyutrvsT+3aSxu26PmtVa7SdlfC5UbTU34cmnAgoRrXVVw5u9U7ttIGn03Xe014n3qGMs4KAK2kugSBpITdaTuIVJxX4mSaGGIRC+rYsYSH0QUI6Pfo4eej+xS6TzvBhGQxgTQ6RBp4W0yD3jkCsZde9H1/gnZKJt1AnPRJD/nDcRm46LwejRhtXpNpYDkv5lOmX2z6dr1FX8JsXLsc7WxMtOiNDkm1ENPFm3kx++dd/+M2Aw81Nk1QXrf0t+Ocn4kueLHRVqDdw0aGmp1HzU36U1gUAczhWJ3HTIInaPxEVTvIPkRyB9+cb6ocYF0jUOwBOoB1PVy8qCQh8LuXZxk8xIkQRJJcuI7ajnIaiDHslOS+s1zbQNKqKCznP/AQGyfLeRevyyOR5En26TwgjeayuKyuXqur53Jx/hsS9HT6YQMAAA==" # forces replacement
      ~ id                   = "cefb6f2293b84c1d6bf041b794cccdb6154fe904" -> (known after apply)
        # (3 unchanged attributes hidden)
    }
$ echo '"H4sIAD8mNWMAA5...hsS9HT6YQMAAA=="' > diff.txt
$ diffdecoding -i diff.txt
```
### Example output:
```
Content-Disposition: attachment; filename="example.com.cfg"
 - path: /etc/nginx/conf.d/example.com.conf
-  defer: true

Content-Type: text/x-shellscript
   ...
    2|      -  
     |2     +  # comment
   ...
```
## Getting help

```sh
## Getting help for related command.
diffdecoding --help
```

## CI/CD

This GitHub repository have a GitHub action that create a release thanks to go-releaser.
