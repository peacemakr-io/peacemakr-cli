<p align="center">
  <br>
    <img src="https://admin.peacemakr.io/p_logo.png" width="150"/>
  <br>
</p>

# Peacemakr E2E-Encryption-as-a-Service CLI
Peacemakr's E2E-Encryption-as-a-Service CLI simplifies your data security with E2E-Encryption service and automated key lifecycle management.

You can easily encrypt your data without worrying about backward compatibility, cross platform portability, or changing security requirements.

Our Zero-Trust capability allows you to customize your security strength to meet the highest standard without having to place your trust in Peacemakr as we don’t have the capacity to get your keys and decrypt your data.

## License

The content of this SDK is open source under [Apache License 2.0](https://github.com/peacemakr-io/peacemakr-cli/blob/master/LICENSE).

## How to Install
```sh
brew tap peacemakr-io/peacemakr
brew install peacemakr

# seeing it work
export PEACEMAKR_APIKEY=your-api-key
echo "hello world" | peacemakr-cli -encrypt | peacemakr-cli -decrypt
# expected output: hello world
```

## Quick start with Dockerhub
 (1) Register an account with https://peacemakr.io

 (2) Login to your admin portal, https://admin.peacemakr.io - grab your API Key

 (3) Encrypt using `peacemakr` + ApiKey,
```sh
$ ciphertext=`echo "hello world" | docker run -e PEACEMAKR_APIKEY=your-api-key -i peacemakr/peacemakr-cli ./peacemakr-cli -encrypt`
$ echo $ciphertext
AQAAAAkAAAAoAAAALAAAADAAAABEAAAAWAAAANcAAADnAAAAMwEAAAAEAAIAAAAAEAAAACEI
eELxb13s32PdZi/4NuUQAAAAUWJ17eT23DpJ63GdnJlq5XsAAAB7ImNyeXB0b0tleUlEIjoi
MGtKdHR6TWt2MnNVa2Q4VndBOXNBckdVaWFZOHI2MHgyV3Y5T29EWGk5QT0iLCJzZW5kZXJL
ZXlJRCI6IjBtaHRscVdMTlM2VUNBWWNRZHk1MG81UTZ5emZSQnltVGZXcHJ3S2JBQ289In0M
AAAAswnudtyjo5K+UOKCSAAAADBGAiEA6b68OXaUWdvyfQrr6jENzjhwn7ewXp9tKNYEmu/W
1rMCIQDJouUC0qlmDhUKpZr1k7gz3zDYuaZsMQs4RH2A2xhadvKfNCD/rk2UG2NVkLQSBHjF
VK2LIDbz40rgi5fdY38C
```
(4) Decrypt using `peacemakr-cli` + ApiKey,
```sh
echo "$ciphertext" | docker run -e PEACEMAKR_APIKEY=your-api-key -i peacemakr/peacemakr-cli ./peacemakr-cli -decrypt 2>/dev/null

hello world
```


Don't want to use docker? Checkout out our native binary releases for ubuntu.
 * https://github.com/peacemakr-io/peacemakr-cli/releases

### Examples in action:
 * [How to encrypt server logs with logrotate and peacemakr-cli](https://medium.com/@danielhuang37/encrypting-all-your-logs-in-2-easy-steps-using-logrotate-and-peacemakr-8ad9cbfe1b4c)

## Build from source with docker
```sh
./build-dep.sh
./build-bin.sh
```

## Setup for Ubuntu
```sh
// get peacemakr cli
curl -LJO https://github.com/peacemakr-io/peacemakr-cli/releases/download/v0.3.0/peacemakr-cli-linux.tar.gz
tar -xf peacemakr-cli-linux.tar.gz
// install dependencies
sudo apt-get install musl musl-dev
ln -s /usr/lib/x86_64-linux-musl/libc.so /usr/lib/libc.musl-x86_64.so.1

// get corecrypt, move it to /usr/lib
curl -LJO https://github.com/peacemakr-io/peacemakr-core-crypto/releases/download/v0.2.2/peacemakr-core-crypto-ubuntu-x86_64.tar.gz
tar -xf peacemakr-core-crypto-ubuntu-x86_64.tar.gz
// move lib/*.so to /usr/lib/
cp -r lib/* /usr/lib/
// move include/* to /usr/include
cp -r include/* /usr/include/

wget https://raw.githubusercontent.com/peacemakr-io/peacemakr-cli/master/peacemakr.yml

export PEACEMAKR_APIKEY=*****client-api-key*****
echo "hello secure world" | ./peacemakr-cli --encrypt
```

## What flags does it accept?
```sh
$ docker run -i peacemakr/peacemakr-cli ./peacemakr -help
Usage of ./peacemakr-cli:
  -config string
        custom config file e.g. (peacemakr.yml) (default "peacemakr.yml")
  -decrypt
        Should the application decrypt the ciphertext
  -domain -domain=DOMAIN_NAME
        A specific use domain to encrypt; -domain=DOMAIN_NAME
  -encrypt
        Should the application encrypt the message
  -inputFileName string
        inputFile to encrypt/decrypt
  -is-peacemakr-blob
        Should the application validate whether the blob is a Peacemakr blob or not
  -outputFileName string
        outputFile to encrypt/decrypt
  -signOnly
        Should the application sign the message
  -verifyOnly
        Should the application verify the input blob
```

## Contributing
We appreciate all contributions. Some basic guidelines are here, for more informaton
see CONTRIBUTING.md

Issues:
- Please include a minimal example that reproduces your issue
- Please use the tags to help us help you
- If you file an issue and you want to work on it, fantastic! Please assign it to yourself.

PRs:
- All PRs must be reviewed and pass CI
- New functionality must include new tests

## FAQ

 * `failed to register due to unknown error (status 500)` - verify correctness of your API key
