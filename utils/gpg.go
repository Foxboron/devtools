package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	openpgp "golang.org/x/crypto/openpgp"
)

// This doesn't work. Need actual public key...
func VerifySignatureOpenPGP(filepath, signaturePath string) {
	signature, err := os.Open(signaturePath)
	if err != nil {
		log.Fatal(err)
	}
	key, err := os.Open(signaturePath)
	if err != nil {
		log.Fatal(err)
	}
	verifyTarget, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	keyIssues, err := openpgp.ReadKeyRing(key)
	if err != nil {
		log.Fatal(err)
	}
	entity, err := openpgp.CheckArmoredDetachedSignature(keyIssues, verifyTarget, signature)
	if err != nil {
		entity, err := openpgp.CheckDetachedSignature(keyIssues, verifyTarget, signature)
		if err != nil {
			fmt.Println("lol")
			fmt.Println(entity)
			log.Fatal()
		}
		fmt.Println(entity.PrimaryKey)
	}
	fmt.Println(entity.PrimaryKey)
}

// Verify signature with gnupg
func VerifySignature(filepath, signaturePath string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("/usr/bin/gpg", "--verify", signaturePath, filepath)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
