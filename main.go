package main

import (
	"crypto"
	"encoding/pem"
	"flag"
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

var (
	printPrivToStdErr bool
	writeToFiles      bool
	forceStupidness   bool
	pubKeyFile        string
	privKeyFile       string
)

func main() {
	flag.BoolVar(&printPrivToStdErr, "privtoerr", false, "When true, the tool prints private key to stderr and public to stdout")
	flag.BoolVar(&writeToFiles, "write", true, "When true it writes the private and public keys to file")
	flag.BoolVar(&forceStupidness, "force", false, "When true, you ignore the author's recommendation about seed size\n(should be at minimum 32 chars, more is better) and continues with your stupidness")
	flag.StringVar(&pubKeyFile, "pubfile", "id_ed25519.pub", "Filename to write public key to")
	flag.StringVar(&privKeyFile, "privfile", "id_ed25519", "Filename to write private key to")
	flag.Parse()
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		flag.Usage()
		fmt.Println("\n\nThe command is intended to work with pipes.")
		fmt.Println("Example: \n\techo 'never lose a key again 81S1r8zpVuFjpJ5odwDTmplp4HZ5JskQ' | anything2ed25519")
		return
	}

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Failed to read from STDIN: %v\n", err)
		return
	}

	if len(bytes) < 32 {
		fmt.Println("WARNING!\nWARNING!\nWARNING!\nWARNING!\nWARNING!\nWARNING!\nWARNING!\nWARNING!\nWARNING!\nWARNING!\nWARNING!\nWARNING!\n")
		fmt.Println("You should use some stronger seed than this! -force to override this idiot protection :)")
	}

	publicKey, privateKey, err := GenerateKeys(bytes)
	if err != nil {
		fmt.Printf("Failed to generate keys: %v\n", err)
		return
	}

	if printPrivToStdErr {
		fmt.Fprintf(os.Stderr, "%s\n", privateKey)
	} else {
		fmt.Printf("%s\n", privateKey)
	}
	fmt.Printf("%s\n", publicKey)

	if writeToFiles {
		_ = ioutil.WriteFile(privKeyFile, privateKey, 0600)
		_ = ioutil.WriteFile(pubKeyFile, publicKey, 0644)
	}
}
