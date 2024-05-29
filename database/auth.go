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

func SignUp(db *sql.DB, email string, username string, password string, secretCode string) (int, int, bool, bool, bool, error) {
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
	var isEmailExist, isUsernameExist, isVerified bool

	sql_statement := `
	CALL signup($1, $2, $3, $4, $5,
	$6, $7, $8, $9, $10);
	`

	err := db.QueryRow(sql_statement, username, email, encryptedPassword, saltString, secretCode, &userId, &emailId, &isEmailExist, &isUsernameExist, &isVerified).
		Scan(&userId, &emailId, &isEmailExist, &isUsernameExist, &isVerified)
	if err != nil {
		panic(err)
	}

	return userId, emailId, isEmailExist, isUsernameExist, isVerified, nil
}

func LogIn(db *sql.DB, nameOrEmail string, password string) (bool, bool, bool, int, string, error) {
	var (
		encryptedPassword string
		saltString        string
		uid               int
		isVerified        bool
		isSignedUp        bool
	)

	sql_statement := `
	SELECT password, salt, id, is_verified 
	FROM Users 
	WHERE LOWER(email) = LOWER($1) OR LOWER(username) = LOWER($1);
	`

	rows, err := db.Query(sql_statement, nameOrEmail)

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

// Change password returning whether the user exist, whether the account has been verified, and whether the old password match
func ChangePassword(db *sql.DB, userid int, oldPassword string, newPassword string) (bool, bool, bool, error) {
	sql_statement := `
	SELECT password, salt, is_verified FROM users WHERE id = $1;
	`

	var encryptedPassword, saltString string
	var isVerified bool
	rows, err := db.Query(sql_statement, userid)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if !rows.Next() {
		return false, false, false, nil
	}

	rows.Scan(&encryptedPassword, &saltString, &isVerified)

	if !isVerified {
		return true, false, false, nil
	}

	salt, err1 := hex.DecodeString(saltString)
	if err1 != nil {
		panic(err1)
	}

	if !doPasswordsMatch(encryptedPassword, oldPassword, []byte(salt)) {
		return true, true, false, nil
	}

	UpdatePassword(db, userid, newPassword)

	return true, true, true, nil
}

func UpdatePassword(db *sql.DB, userid int, newPassword string) (bool, error) {
	sql_statement := `
	UPDATE users
	SET password = $1, salt = $2
	WHERE id = $3;
	`
	newSalt := generateRandomSalt()
	newSaltString := hex.EncodeToString(newSalt)
	newEncryptedPassword := hashPassword(newPassword, []byte(newSalt))

	_, err := db.Exec(sql_statement, newEncryptedPassword, newSaltString, userid)
	if err != nil {
		panic(err)
	}
	return true, nil
}

// Make password recovery returning whether the user exist, has been verified, the recovery id, and the user's email
func MakePasswordRecovery(db *sql.DB, email string, secretCode string) (bool, bool, int, string, error) {
	sql_get_user_id := `
	SELECT u.id
	FROM users u
	WHERE LOWER(u.email) = LOWER($1);
	`

	var userid int
	err := db.QueryRow(sql_get_user_id, email).Scan(&userid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, false, -1, "", errors.New("email does not exist")
		}
		return false, false, -1, "", errors.New("something went wrong")
	}

	sql_statement := `
	CALL make_password_recovery($1, $2, $3, $4, $5, $6);
	`

	var isExist, isVerified bool
	var recoveryId int
	var email2 string
	err2 := db.QueryRow(sql_statement, userid, secretCode, &isExist, &isVerified, &recoveryId, &email2).Scan(&isExist, &isVerified, &recoveryId, &email2)
	if err2 != nil {
		return false, false, -1, "", errors.New("something went wrong")
	}
	return isExist, isVerified, recoveryId, email2, nil
}

// Get whether the recoveryId and secretCode match, expired, and used
func GetRecoverPassword(db *sql.DB, recoveryId int, secretCode string) (bool, bool, bool, error) {
	sql_statement := `
	SELECT secret_code = $1 AS is_match, expired_at :: time < CURRENT_TIME AS is_expired, is_used
	FROM password_recoveries
	WHERE id = $2;
	`

	var isMatch, isExpired, isUsed bool
	err := db.QueryRow(sql_statement, secretCode, recoveryId).Scan(&isMatch, &isExpired, &isUsed)
	if err != nil {
		panic(err)
	}

	return isMatch, isExpired, isUsed, nil
}

// Recover the password return whether the password recovery exist, expired, match, or used
func RecoverPassword(db *sql.DB, recoveryId int, secretCode string, newPassword string) (bool, bool, bool, bool, error) {
	sql_statement := `
	CALL recover_password($1, $2, $3, $4, $5, $6, $7, $8);
	`

	var isExist, isExpired, isMatch, used bool

	newSalt := generateRandomSalt()
	newSaltString := hex.EncodeToString(newSalt)
	newEncryptedPassword := hashPassword(newPassword, []byte(newSalt))

	err := db.QueryRow(sql_statement, recoveryId, secretCode, newEncryptedPassword, newSaltString, &isExist, &isExpired, &isMatch, &used).
		Scan(&isExist, &isExpired, &isMatch, &used)
	if err != nil {
		panic(err)
	}

	return isExist, isExpired, isMatch, used, nil
}
