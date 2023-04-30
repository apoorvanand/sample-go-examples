package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

func main() {
    // Generate a 32-byte random string for the public key
    publicKey := make([]byte, 32)
    if _, err := rand.Read(publicKey); err != nil {
        panic(err)
    }

    // Generate a 32-byte random string for the private key
    privateKey := make([]byte, 32)
    if _, err := rand.Read(privateKey); err != nil {
        panic(err)
    }

    // Encode the keys to base64 strings
    publicKeyStr := base64.URLEncoding.EncodeToString(publicKey)
    privateKeyStr := base64.URLEncoding.EncodeToString(privateKey)

    // Print the keys
    fmt.Printf("Public key: %s\n", publicKeyStr)
    fmt.Printf("Private key: %s\n", privateKeyStr)
}
