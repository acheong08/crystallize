# crystallize
Encryption and signing tool using Kyber and Dilithium

# Usage
```
crystallize v0.0.1 - A CLI tool for post quantum cryptography

Available commands:

   generate     Generate a keypair 
   encryption   Encrypt or decrypt data 
   signing      Sign or verify data 

Flags:

  -help
    	Get help on the 'crystallize' command.
```
```
crystallize generate - Generate a keypair
Flags:

  -dili
    	Use Dilithium
  -help
    	Get help on the 'crystallize generate' command.
  -kyber
    	Use Kyber
  -privkey string
    	The private key file (default "priv.key")
  -pubkey string
    	The public key file (default "pub.key")
```
```
crystallize encryption - Encrypt or decrypt data
Flags:

  -decrypt
    	Decrypt the data
  -encrypt
    	Encrypt the data
  -file string
    	The file to process
  -help
    	Get help on the 'crystallize encryption' command.
  -output string
    	The output file
  -privkey string
    	The private key to use (default "priv.key")
  -pubkey string
    	The public to use (default "pub.key")
  -text string
    	The text to process
```
```
crystallize signing - Sign or verify data
Flags:

  -file string
    	The file to process
  -help
    	Get help on the 'crystallize signing' command.
  -output string
    	The output file
  -privkey string
    	The private key to use (default "priv.key")
  -pubkey string
    	The public to use (default "pub.key")
  -sign
    	Sign the data
  -signature string
    	The signature file to verify
  -text string
    	The text to process
  -verify
    	Verify the data
```

# Installation
`go install github.com/acheong08/crystallize@latest` 
