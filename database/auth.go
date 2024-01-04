package database

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"example/web-service-gin/token"
	"fmt"
	"time"

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

func SignUp(db *sql.DB, email string, username string, password string) (int, string, error) {
	salt := generateRandomSalt()
	saltString := hex.EncodeToString(salt)
	encryptedPassword := hashPassword(password, salt)

	nexId := getUserId(db)
	_, err := db.Exec("INSERT INTO Users (id, email, username, password, salt) VALUES ($1, $2, $3, $4, $5)", nexId, email, username, encryptedPassword, saltString)
	if err != nil {
		return -1, "", err
	}

	token, err2 := token.GenerateToken(nexId)
	if err2 != nil {
		return -1, "", err2
	}

	return nexId, token, nil
}

func LogIn(db *sql.DB, nameOrEmail string, password string) (bool, int, string, error) {
	var (
		encryptedPassword string
		saltString        string
		uid               int
	)

	rows, err := db.Query("SELECT password, salt, id FROM Users WHERE email = $1 OR username = $1", nameOrEmail)
	if rows == nil {
		return false, -1, "wrong username or password", nil
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&encryptedPassword, &saltString, &uid)
		break
	}

	if err != nil {
		fmt.Println("Err here", err.Error())
		panic(err)
	}

	salt, err2 := hex.DecodeString(saltString)
	if err2 != nil {
		panic(err2)
	}

	success := doPasswordsMatch(encryptedPassword, password, salt)
	if !success {
		return success, -1, "", nil
	}

	token, err3 := token.GenerateToken(uid)
	if err3 != nil {
		panic(err3)
	}

	userid, err := getUserIdFromNameOrEmail(db, nameOrEmail)
	if err != nil {
		panic(err)
	}

	err = storeUserIdAndJWT(db, userid, token)
	if err != nil {
		panic(err)
	}

	return success, userid, token, nil
}

func checkToken(db *sql.DB, userId int, token string) error {
	sql_statement := `
	SELECT token
	FROM tokens
	WHERE userid = $1 AND token = $2
	`

	rows, err := db.Query(sql_statement, userId, token)

	if err != nil {
		return errors.New("Unauthorized")
	}

	isUserFound := false

	for rows.Next() {
		isUserFound = true
		break
	}

	if !isUserFound {
		return errors.New("Unauthorized")
	}

	return nil
}

func storeUserIdAndJWT(db *sql.DB, userid int, token string) error {
	sql_statement := `
	INSERT INTO tokens(userid, token) VALUES($1, $2) ON CONFLICT DO NOTHING;
	`
	// fmt.Println(time.Now().Format("2025-01-02"))
	rows, err := db.Query(sql_statement, userid, token)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return nil
}

// Store email verification secret code returning the row id
func StoreSecretCode(db *sql.DB, userid int, email string, secretCode string) (int, error) {
	sql_statement := `
	INSERT INTO Email_verifications (id, user_id, email, secret_code) 
	VALUES ((SELECT COUNT(*) FROM Email_verifications) + 1, $1, $2, $3)
	RETURNING id;
	`

	var id int

	err := db.QueryRow(sql_statement, userid, email, secretCode).Scan(&id)

	if err != nil {
		panic(err)
	}
	return id, nil
}

// Get email verification's secret code and whether it is expired
func GetSecretCode(db *sql.DB, emailid int) (string, bool, error) {
	sql_statement := `
	SELECT secret_code, expired_at
	FROM email_verifications
	WHERE id = $1;
	`

	var secretCode string
	var expiredAt time.Time

	err := db.QueryRow(sql_statement, emailid).Scan(&secretCode, &expiredAt)
	if err != nil {
		panic(err)
	}

	return secretCode, time.Now().After(expiredAt), nil
}

// Verify email if valid returning whether the email already verified, secret code expired, and secret code does not match
func VerifyEmail(db *sql.DB, emailid int, secretCode string) (bool, bool, bool, error) {
	sql_statement := `
	CALL verify_email($1, $2, $3, $4, $5);
	`

	var isVerified, isExpired, isMatch bool
	err := db.QueryRow(sql_statement, emailid, secretCode, &isVerified, &isExpired, &isMatch).Scan(&isVerified, &isExpired, &isMatch)
	if err != nil {
		panic(err)
	}

	return isVerified, isExpired, isMatch, nil
}
