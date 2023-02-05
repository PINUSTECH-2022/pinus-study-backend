package database

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"example/web-service-gin/token"

	_ "github.com/lib/pq"
)

// Generate 16 bytes randomly and securely using the
// Cryptographically secure pseudorandom number generator (CSPRNG)
// in the crypto.rand package
func generateRandomSalt() []byte {
	var salt = make([]byte, 16)

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
func doPasswordsMatch(hashedPassword, currPassword string, salt []byte) bool {
	var currPasswordHash = hashPassword(currPassword, salt)

	return hashedPassword == currPasswordHash
}

func SignUp(db *sql.DB, email string, username string, password string) error {
	salt := generateRandomSalt()
	saltString := hex.EncodeToString(salt)
	encryptedPassword := hashPassword(password, salt)

	rows, err := db.Query("INSERT INTO Users (id, email, username, password, salt) VALUES ($1, $2, $3, $4, $5)", getUserId(db), email, username, encryptedPassword, saltString)
	if err != nil {
		return errors.New(err.Error() + "12312312")
	}
	defer rows.Close()
	return nil
}

func LogIn(db *sql.DB, nameOrEmail string, password string) (bool, string, error) {
	rows, err := db.Query("SELECT password FROM Users WHERE email = $1 OR username = $1", nameOrEmail)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var encryptedPassword string

	for rows.Next() {
		err := rows.Scan(&encryptedPassword)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	rows, err = db.Query("SELECT salt FROM Users WHERE email = $1 OR username = $1", nameOrEmail)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var saltString string

	for rows.Next() {
		err := rows.Scan(&saltString)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	salt, err := hex.DecodeString(saltString)
	if err != nil {
		panic(err)
	}

	success := doPasswordsMatch(encryptedPassword, password, salt)

	rows, err = db.Query("SELECT id FROM Users WHERE email = $1 OR username = $1", nameOrEmail)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var uid int

	for rows.Next() {
		err := rows.Scan(&uid)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	token, err := token.GenerateToken(uid)
	if err != nil {
		panic(err)
	}

	return success, token, nil
}
