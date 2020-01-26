<p align="center">
  <br>
    <img src="https://github.com/peacemakr-io/peacemakr-admin-portal/blob/master/peacemakr-admin/public/images/PeacemakrP-Golden.png" width="150"/>
  <br>
</p>

# peacemakr-cli
CLI that may be used to interact with anything protected by Peacemakr's Secure Data Platform (includes builtin Key Lifecycle Management, Forward Cryptoagility, and On-Prem Key Derivation).

## Quick start with Dockerhub
 (1) Register an account with https://peacemakr.io
 
 (2) Login to your admin portal, https://admin.peacemakr.io - grab your ApiKey

 (3) Encrypt using `peacemakr-cli` + ApiKey,
```sh
$ ciphertext=`echo "hello world" | docker run -e PEACEMAKR_APIKEY=your-api-key -i peacemakr/peacemakr-cli ./peacemakr-cli -action=encrypt`
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
echo "$ciphertext" | docker run -e PEACEMAKR_APIKEY=your-api-key -i peacemakr/peacemakr-cli ./peacemakr-cli -action=decrypt 2>/dev/null

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

## What flags does it accept?
```sh
$ docker run -i peacemakr/peacemakr-cli ./peacemakr-cli -help
Usage of ./peacemakr-cli:
  -action string
    	action= encrypt|decrypt (default "encrypt")
  -config string
    	custom config file e.g. (peacemakr.yml) (default "peacemakr.yml")
  -inputFileName string
    	inputFile to encrypt/decrypt
  -outputFileName string
    	outputFile to encrypt/decrypt
```

## FAQ

 * `failed to register due to unknown error (status 500)` - verify correctness of your API key
