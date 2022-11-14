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

	// Flag vars
	var text string
	var file string
	var output string
	var priv string = "priv.key"
	var pub string = "pub.key"
	var encryptFlag bool
	var decryptFlag bool
	var signFlag bool
	var verifyFlag bool
	var signature string
	var dilithiumFlag bool
	var kyberFlag bool

	/// Generate a keypair
	generateCmd := cli.NewSubCommand("generate", "Generate a keypair")

	generateCmd.StringFlag("pubkey", "The public key file", &pub)
	generateCmd.StringFlag("privkey", "The private key file", &priv)
	generateCmd.BoolFlag("dili", "Use Dilithium", &dilithiumFlag)
	generateCmd.BoolFlag("kyber", "Use Kyber", &kyberFlag)

	generateCmd.Action(func() error {
		// Checking for flags
		if dilithiumFlag && kyberFlag {
			return fmt.Errorf("cannot use both Dilithium and Kyber")
		} else if !dilithiumFlag && !kyberFlag {
			return fmt.Errorf("must specify either Dilithium or Kyber")
		}
		if dilithiumFlag {
			pk, sk := security.GenerateDili()
			script.Echo(encode(pk)).WriteFile(pub)
			script.Echo(encode(sk)).WriteFile(priv)
		} else if kyberFlag {
			pk, sk := security.GenerateKyber()
			script.Echo(string(pk)).WriteFile(pub)
			script.Echo(string(sk)).WriteFile(priv)
		}
		fmt.Println("Generated keypair")
		return nil
	})

	/// For encryption
	encryptionCmd := cli.NewSubCommand("encryption", "Encrypt or decrypt data")
	// Flags
	encryptionCmd.StringFlag("text", "The text to process", &text)
	encryptionCmd.StringFlag("file", "The file to process", &file)
	encryptionCmd.StringFlag("output", "The output file", &output)
	encryptionCmd.StringFlag("privkey", "The private key to use", &priv)
	encryptionCmd.StringFlag("pubkey", "The public to use", &pub)
	encryptionCmd.BoolFlag("encrypt", "Encrypt the data", &encryptFlag)
	encryptionCmd.BoolFlag("decrypt", "Decrypt the data", &decryptFlag)

	encryptionCmd.Action(func() error {
		// Checking for flags
		if encryptFlag && decryptFlag {
			return fmt.Errorf("cannot encrypt and decrypt at the same time")
		} else if !encryptFlag && !decryptFlag {
			return fmt.Errorf("must specify either encrypt or decrypt")
		}
		if text == "" && file == "" {
			return fmt.Errorf("must specify either text or file")
		}
		if text != "" && file != "" {
			return fmt.Errorf("cannot specify both text and file")
		}

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
			fmt.Println("Wrote to encrypted file")
		} else {
			fmt.Println(out)
		}
		return nil
	})

	/// For signing
	signCmd := cli.NewSubCommand("signing", "Sign or verify data")

	// Flags
	signCmd.StringFlag("text", "The text to process", &text)
	signCmd.StringFlag("file", "The file to process", &file)
	signCmd.StringFlag("output", "The output file", &output)
	signCmd.StringFlag("privkey", "The private key to use", &priv)
	signCmd.StringFlag("pubkey", "The public to use", &pub)
	signCmd.BoolFlag("sign", "Sign the data", &signFlag)
	signCmd.BoolFlag("verify", "Verify the data", &verifyFlag)
	signCmd.StringFlag("signature", "The signature file to verify", &signature)

	signCmd.Action(func() error {
		// Checking for flags
		if signFlag && verifyFlag {
			return fmt.Errorf("cannot sign and verify at the same time")
		} else if !signFlag && !verifyFlag {
			return fmt.Errorf("must specify either sign or verify")
		}
		if text == "" && file == "" {
			return fmt.Errorf("must specify either text or file")
		}
		if text != "" && file != "" {
			return fmt.Errorf("cannot specify both text and file")
		}
		var out []byte
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
		if signFlag {
			//Get private key
			privKey, err := readFromFile(priv)
			privKey = decode(string(privKey))
			if err != nil {
				fmt.Println("Error: Could not read private key")
				return nil
			}
			out = security.Sign(privKey, []byte(data))
		}
		if verifyFlag {
			// Check for signature
			if signature == "" {
				return fmt.Errorf("must specify signature file")
			}
			// Get public key
			pubKey, err := readFromFile(pub)
			pubKey = decode(string(pubKey))
			if err != nil {
				fmt.Println("Error: Could not read public key")
				return nil
			}
			signature_bytes, err := readFromFile(signature)
			if err != nil {
				fmt.Println("Error: Could not read signature")
				return nil
			}
			if security.Verify(pubKey, []byte(data), signature_bytes) {
				fmt.Println("Signature is valid")
			} else {
				fmt.Println("Signature is invalid")
			}
			out = []byte("Signature verified")
		}
		if output != "" {
			// Write to file
			script.Echo(string(out)).WriteFile(output)
		} else {
			fmt.Println(string(out))
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
