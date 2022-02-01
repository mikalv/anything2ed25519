package main

import (
	"crypto"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	// This is just a forked golang.org/x/crypto/ed25519/internal/edwards25519 repo (find it at https://go.googlesource.com/crypto - ed25519/internal)
	// so that we can use the internal API.
	"0xcc.re/anything2ed25519/edwards25519"

	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

func MakeHash(data []byte) []byte {
	var hash crypto.Hash
	hash = crypto.SHA256
	h := hash.New()
	h.Write(data)
	digest := h.Sum(nil)
	return digest
}

func GenerateKeys(bytes []byte) ([]byte, []byte, error) {
	// Generate a new private/public keypair for OpenSSH
	privKey := GenPrivKeyFromSecret(bytes)
	pubKey := getPublicKey(privKey)
	publicKey, err := ssh.NewPublicKey(pubKey)

	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return nil, nil, err
	}

	pemKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privKey),
	}
	privateKey := pem.EncodeToMemory(pemKey)
	formattedPublicKey := ssh.MarshalAuthorizedKey(publicKey)

	fmt.Printf("%s\n", privateKey)
	fmt.Printf("%s\n", formattedPublicKey)
	return formattedPublicKey, privateKey, nil
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
	seed := MakeHash(secret)
	privKey := ed25519.NewKeyFromSeed(seed)
	return privKey
}

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: echo 'never lose a key again 81S1r8zpVuFjpJ5odwDTmplp4HZ5JskQ' | anything2ed25519")
		return
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Failed to read from STDIN: %v\n", err)
		return
	}

	publicKey, privateKey, err := GenerateKeys(bytes)
	if err != nil {
		fmt.Printf("Failed to generate keys: %v\n", err)
		return
	}

	_ = ioutil.WriteFile("id_ed25519", privateKey, 0600)
	_ = ioutil.WriteFile("id_ed25519.pub", publicKey, 0644)
}
