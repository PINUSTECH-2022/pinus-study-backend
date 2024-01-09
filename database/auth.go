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

func SignUp(db *sql.DB, email string, username string, password string, secretCode string) (int, int, bool, bool, error) {
	salt := generateRandomSalt()
	saltString := hex.EncodeToString(salt)
	encryptedPassword := hashPassword(password, salt)

	// nexId := getUserId(db)
	// _, err := db.Exec("INSERT INTO Users (id, email, username, password, salt) VALUES ($1, $2, $3, $4, $5)", nexId, email, username, encryptedPassword, saltString)
	// if err != nil {
	// 	return -1, err
	// }

	// return nexId, nil

	var userId, emailId int
	var isEmailExist, isUsernameExist bool

	sql_statement := `
	CALL signup($1, $2, $3, $4, $5,
	$6, $7, $8, $9);
	`

	err := db.QueryRow(sql_statement, username, email, encryptedPassword, saltString, secretCode, &userId, &emailId, &isEmailExist, &isUsernameExist).
		Scan(&userId, &emailId, &isEmailExist, &isUsernameExist)
	if err != nil {
		panic(err)
	}

	return userId, emailId, isEmailExist, isUsernameExist, nil
}

func LogIn(db *sql.DB, nameOrEmail string, password string) (bool, bool, bool, int, string, error) {
	var (
		encryptedPassword string
		saltString        string
		uid               int
		isVerified        bool
		isSignedUp        bool
	)

	rows, err := db.Query("SELECT password, salt, id, is_verified FROM Users WHERE email = $1 OR username = $1", nameOrEmail)

	defer rows.Close()

	if !rows.Next() {
		isSignedUp = false
		return false, isSignedUp, false, -1, "", nil
	} else {
		isSignedUp = true
		rows.Scan(&encryptedPassword, &saltString, &uid, &isVerified)
	}

	if err != nil {
		fmt.Println("Err here", err.Error())
		panic(err)
	}

	// Check whether the email has been verified
	if !isVerified {
		return false, isSignedUp, isVerified, -1, "", nil
	}

	salt, err2 := hex.DecodeString(saltString)
	if err2 != nil {
		panic(err2)
	}

	success := doPasswordsMatch(encryptedPassword, password, salt)
	if !success {
		return success, isSignedUp, isVerified, -1, "", nil
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

	return success, isSignedUp, isVerified, userid, token, nil
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

// Stores email verification secret code returning whether the userid exist, the email id, the email, the username, and whether the account has been verified
func StoreSecretCode(db *sql.DB, userid int, secretCode string) (bool, int, string, string, bool, error) {
	sql_statement := `
	CALL make_verification($1, $2, 
	$3, $4, $5, $6, $7);
	`

	var id int
	var email, username string
	var isExist, isVerified bool

	err := db.QueryRow(sql_statement, userid, secretCode, &isExist, &id, &email, &username, &isVerified).
		Scan(&isExist, &id, &email, &username, &isVerified)

	if err != nil {
		panic(err)
	}
	return isExist, id, email, username, isVerified, nil
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
