package main

import (
	"encoding/base64"
	"fmt"

	security "github.com/acheong08/crystallize/security"
	script "github.com/bitfield/script"
	clir "github.com/leaanthony/clir"
)

func main() {
	// Create a new CLI
	cli := clir.NewCli("crystallize", "A CLI tool for post quantum cryptography", "v0.0.1")

	// Add subcommands
	generateCmd := cli.NewSubCommand("generate", "Generate a keypair")

	var pubkey_path string
	var privkey_path string
	generateCmd.StringFlag("pubkey", "The public key file", &pubkey_path)
	generateCmd.StringFlag("privkey", "The private key file", &privkey_path)

	generateCmd.Action(func() error {
		pk, sk := security.Generate()
		script.Echo(string(pk)).WriteFile(pubkey_path)
		script.Echo(string(sk)).WriteFile(privkey_path)
		return nil
	})

	// Flags
	var text string
	cli.StringFlag("text", "The text to process", &text)
	var file string
	cli.StringFlag("file", "The file to process", &file)
	var output string
	cli.StringFlag("output", "The output file", &output)
	var priv string
	cli.StringFlag("privkey", "The private key to use", &priv)
	var pub string
	cli.StringFlag("pubkey", "The public to use", &pub)
	var encryptFlag bool
	cli.BoolFlag("encrypt", "Encrypt the data", &encryptFlag)
	var decryptFlag bool
	cli.BoolFlag("decrypt", "Decrypt the data", &decryptFlag)

	// Main actions
	cli.Action(func() error {
		var out string
		var data string
		if text != "" {
			data = encode([]byte(text))
		} else if file != "" {
			data_bytes, err := readFromFile(file)
			data = encode(data_bytes)
			if err != nil {
				return err
			}
		}
		if encryptFlag {
			// Get public key
			pubKey, err := readFromFile(pub)
			if err != nil {
				fmt.Println("Error: Could not read public key")
				return nil
			}
			out = security.Encrypt(pubKey, string(data))
		}
		if decryptFlag {
			// Get private key
			privKey, err := readFromFile(priv)
			if err != nil {
				fmt.Println("Error: Could not read private key")
				return nil
			}
			data = string(decode(data))
			out = string(decode(security.Decrypt(privKey, data)))
		}
		if output != "" {
			// Write to file
			script.Echo(out).WriteFile(output)
		}
		return nil
	})

	// Run the CLI
	err := cli.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func readFromFile(file string) ([]byte, error) {
	return script.File(file).Bytes()
}

func encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func decode(data string) []byte {
	decoded, _ := base64.StdEncoding.DecodeString(data)
	return decoded
}
