# Make a SSH Ed25519 private/public keypair of **anything** you can pipe into a UNIX pipe

## Compile

```
mkdir go
export GOLANG=$(pwd)/go
cd go ; mkdir -p src/0xcc.re
cd src/0xcc.re
git clone https://github.com/mikalv/anything2ed25519.git
cd anything2ed25519
./build.sh
```

## Prebuilt binaries

Check out the directory `bin`

## Notes

The crypto is safe, however don't be stupid. echo'ing "password", or "abc" or "hello world" or anything common into anything2ed25519 - **don't ever** use that anywhere, since people can guess your private key.. This is like a bitcoin "mnemonic code" and should have at least 24 words..


## Author

Mikal Villa <mikalv@mikalv.net>

