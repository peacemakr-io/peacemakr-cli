# peacemakr-cli
CLI that encrypts and decrypts anything using Peacemakr Crypto System.

## Quick start
 (1) Register an account with https://peacemakr.io
 (2) Login to your admin portal, https://admin.peacemakr.io
 (3) Grab your ApiKey
 (4) Invoke `peacemakr-cli` using your ApiKey, for example,
```sh
$ echo "hello world" | docker run -e PEACEMAKR_APIKEY=your-api-key -i peacemakr-cli ./peacemakr-cli
AQAAAAkAAAAoAAAALAAAADAAAABEAAAAWAAAANcAAADnAAAAMwEAAAAEAAIAAAAAEAAAACEI
eELxb13s32PdZi/4NuUQAAAAUWJ17eT23DpJ63GdnJlq5XsAAAB7ImNyeXB0b0tleUlEIjoi
MGtKdHR6TWt2MnNVa2Q4VndBOXNBckdVaWFZOHI2MHgyV3Y5T29EWGk5QT0iLCJzZW5kZXJL
ZXlJRCI6IjBtaHRscVdMTlM2VUNBWWNRZHk1MG81UTZ5emZSQnltVGZXcHJ3S2JBQ289In0M
AAAAswnudtyjo5K+UOKCSAAAADBGAiEA6b68OXaUWdvyfQrr6jENzjhwn7ewXp9tKNYEmu/W
1rMCIQDJouUC0qlmDhUKpZr1k7gz3zDYuaZsMQs4RH2A2xhadvKfNCD/rk2UG2NVkLQSBHjF
VK2LIDbz40rgi5fdY38C
```
Make it easy with an alias:
```sh
...
```


## Build with docker
```sh
# build dependencies
./build-dep.sh
# build binaries 
./build-bin.sh
```

## Run with docker
```sh
$ docker run -e PEACEMAKR_APIKEY=your-api-key -i peacemakr-cli ./peacemakr-cli -help
```

