// package main

// import "groupietracker.com/m/pkg/server"

// func main() {
// 	server.StartServer()
// }

package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

// Define salt size
const saltSize = 16

// Generate 16 bytes randomly and securely using the
// Cryptographically secure pseudorandom number generator (CSPRNG)
// in the crypto.rand package
func generateRandomSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return salt
}

// Combine password and salt then hash them using the SHA-512
// hashing algorithm and then return the hashed password
// as a hex string
func hashPassword(password string, salt []byte) string {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Create sha-512 hasher
	var sha512Hasher = sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	// Convert the hashed password to a hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

// Check if two passwords match
func doPasswordsMatch(hashedPassword, currPassword string,
	salt []byte) bool {
	var currPasswordHash = hashPassword(currPassword, salt)

	return hashedPassword == currPasswordHash
}

func main() {
	// First generate random 16 byte salt
	var salt = generateRandomSalt(saltSize)

	// Hash password using the salt
	var hashedPassword = hashPassword("hello", salt)

	fmt.Println("Password Hash:", hashedPassword)
	fmt.Println("Salt:", salt)

	// Check if passed password matches the original password by hashing it
	// with the original password's salt and check if the hashes match
	fmt.Println("Password Match:",
		doPasswordsMatch(hashedPassword, "hello", salt))
}