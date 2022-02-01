# Make a SSH Ed25519 private/public keypair of **anything** you can pipe into a UNIX pipe

```
Usage of anything2ed25519:
  -privfile string
    	Filename to write private key to (default "id_ed25519")
  -privtoerr
    	When true, the tool prints private key to stderr and public to stdout
  -pubfile string
    	Filename to write public key to (default "id_ed25519.pub")
  -write
    	When true it writes the private and public keys to file (default true)
```

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

## Example usage

**NOTE** Please don't be stupid and use this key, or anything else with a **WEAK** seed. (As "Hello World" is...)

```
anything2ed25519 [main●] % echo "Hello World" | ./bin/anything2ed25519-darwin-amd64.macho
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtz
c2gtZWQyNTUxOQAAACBkj/xWESur4UT48VFau6Vkk56Kzt48qvJk/UGR2DS7dAAA
AIiaywRCmssEQgAAAAtzc2gtZWQyNTUxOQAAACBkj/xWESur4UT48VFau6Vkk56K
zt48qvJk/UGR2DS7dAAAAEBhNTkxYTZkNDBiZjQyMDQwNGEwMTE3MzNjZmI3YjE5
MGSP/FYRK6vhRPjxUVq7pWSTnorO3jyq8mT9QZHYNLt0AAAAAAECAwQF
-----END OPENSSH PRIVATE KEY-----

ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICv3hIFN8Aw2k16freruJ1Ve4L0V7Z6bKwnYrV7R071+
```


## Prebuilt binaries

Check out the directory `bin`

## Notes

The crypto is experimental, however don't be stupid. echo'ing "password", or "abc" or "hello world" or anything common into anything2ed25519 - **don't ever** use that anywhere, since people can guess your private key.. This is like a bitcoin "mnemonic code" and should have at least 24 words..

## Changelog

* 2022-02-01 patch#2 : Added flags and more usability stuff
* 2022-02-01 patch#1 : Some bugs found and fixed after public release, thanks to people at lobste.rs
* ~3-4yrs ago: Initial write

## Author

Mikal Villa <mikalv@mikalv.net>

