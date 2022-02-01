package main

import (
	"crypto"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	// This is just a forked golang.org/x/crypto/ed25519/internal/edwards25519 repo (find it at https://go.googlesource.com/crypto - ed25519/internal)
	// so that we can use the internal API.
	"0xcc.re/anything2ed25519/edwards25519"

	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

func MakeHash(data []byte) string {
	var hash crypto.Hash
	hash = crypto.SHA256
	h := hash.New()
	h.Write(data)
	digest := h.Sum(nil)
	sh := string(fmt.Sprintf("%x\n", digest))
	return sh
}

// Generate the public key corresponding to the already hashed private
// key.
//
// This code is mostly copied from GenerateKey in the
// golang.org/x/crypto/ed25519 package, from after the SHA512
// calculation of the seed.
func getPublicKey(privateKey []byte) ed25519.PublicKey {
	var A edwards25519.ExtendedGroupElement
	var hBytes [32]byte
	copy(hBytes[:], privateKey)
	edwards25519.GeScalarMultBase(&A, &hBytes)
	var publicKeyBytes [32]byte
	A.ToBytes(&publicKeyBytes)

	return ed25519.PublicKey(publicKeyBytes[:])
}

// GenPrivKeyFromSecret hashes the secret with SHA2, and uses
// that 32 byte output to create the private key.
// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GenPrivKeyFromSecret(secret []byte) ed25519.PrivateKey {
	seed := MakeHash(secret) // Not Ripemd160 because we want 32 bytes.

	privKey := ed25519.NewKeyFromSeed([]byte(seed)[:32])
	return privKey
}

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: echo 'ALDRI mist keys igjen' | ed25519cow")
		return
	}

	//reader := bufio.NewReader(os.Stdin)
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Printf("Failed to read: %v", err)
		return
	}
	// Generate a new private/public keypair for OpenSSH
	privKey := GenPrivKeyFromSecret(bytes)
	pubKey := getPublicKey(privKey)
	//pubKey, privKey, _ := ed25519.GenerateKey(reader)
	publicKey, err := ssh.NewPublicKey(pubKey)

	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return
	}

	pemKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privKey),
	}
	privateKey := pem.EncodeToMemory(pemKey)
	authorizedKey := ssh.MarshalAuthorizedKey(publicKey)

	fmt.Printf("%s\n", privateKey)
	fmt.Printf("%s\n", authorizedKey)

	_ = ioutil.WriteFile("id_ed25519", privateKey, 0600)
	_ = ioutil.WriteFile("id_ed25519.pub", authorizedKey, 0644)
}
