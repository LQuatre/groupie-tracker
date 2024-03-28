package passwordManager

import (
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// Combine password and salt then hash them using the SHA-512
// hashing algorithm and then return the hashed password
// as a hex string
// Hash password
func HashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)
  
	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
	  GenerateFromPassword(passwordBytes, bcrypt.MinCost)
  
	return string(hashedPasswordBytes), err
  }
func DoPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func SaveSaltInFile(salt []byte, filePath string) {
	err := ioutil.WriteFile(filePath , salt, 0644)
	if err != nil {
		log.Fatalf("Failed to save salt in file: %v", err)
	}
}

func ReadSaltFromFile(filePath string) []byte {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
	}

	salt, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read salt from file: %v", err)
	}
	
	return salt
}