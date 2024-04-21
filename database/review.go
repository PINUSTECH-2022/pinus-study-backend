package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type Review struct {
	ModuleId      string
	UserId        int
	Username      string
	Timestamp     string
	Workload      int
	ExpectedGrade string
	ActualGrade   string
	Difficulty    int
	SemesterTaken string
	Lecturer      string
	Content       string
	Suggestion    string
}

func GetReviewByModule(db *sql.DB, moduleId string) []Review {
	rows, err := db.Query(`SELECT R.moduleId, R.userId, R.timestamp, R.workload, 
		R.expectedGrade, R.actualGrade, R.difficulty, R.semesterTaken, R.lecturer, 
		R.content, R.suggestion, U.username
		FROM Reviews AS R, Users AS U
		WHERE R.moduleId = $1 AND R.userId = U.id AND R.is_deleted = false`,
		moduleId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		err := rows.Scan(&review.ModuleId, &review.UserId, &review.Timestamp, &review.Workload,
			&review.ExpectedGrade, &review.ActualGrade, &review.Difficulty, &review.SemesterTaken,
			&review.Lecturer, &review.Content, &review.Suggestion, &review.Username)
		if err != nil {
			panic(err)
		}
		reviews = append(reviews, review)
	}
	if rows.Err() != nil {
		panic(err)
	}
	return reviews
}

func GetReviewByUser(db *sql.DB, userId int) []Review {
	rows, err := db.Query(`SELECT R.moduleId, R.userId, R.timestamp, R.workload, 
		R.expectedGrade, R.actualGrade, R.difficulty, R.semesterTaken, R.lecturer, 
		R.content, R.suggestion, U.username
		FROM Reviews AS R, Users AS U
		WHERE R.userId = $1 AND R.userId = U.id AND R.is_deleted = false`,
		userId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		err := rows.Scan(&review.ModuleId, &review.UserId, &review.Timestamp, &review.Workload,
			&review.ExpectedGrade, &review.ActualGrade, &review.Difficulty, &review.SemesterTaken,
			&review.Lecturer, &review.Content, &review.Suggestion, &review.Username)
		if err != nil {
			panic(err)
		}
		reviews = append(reviews, review)
	}
	if rows.Err() != nil {
		panic(err)
	}
	return reviews
}

func GetReviewByModuleAndUser(db *sql.DB, moduleId string, userId int) (bool, Review) {
	rows, err := db.Query(`SELECT R.moduleId, R.userId, R.timestamp, R.workload, 
		R.expectedGrade, R.actualGrade, R.difficulty, R.semesterTaken, R.lecturer, 
		R.content, R.suggestion, U.username
		FROM Reviews AS R, Users AS U
		WHERE R.moduleId = $1 AND U.id = $2 AND R.userId = U.id AND R.is_deleted = false`,
		moduleId, userId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var review Review
	isReviewExist := false
	for rows.Next() {
		isReviewExist = true
		err := rows.Scan(&review.ModuleId, &review.UserId, &review.Timestamp, &review.Workload,
			&review.ExpectedGrade, &review.ActualGrade, &review.Difficulty, &review.SemesterTaken,
			&review.Lecturer, &review.Content, &review.Suggestion, &review.Username)
		if err != nil {
			panic(err)
		}
	}
	if rows.Err() != nil {
		panic(err)
	}
	return isReviewExist, review
}

func GetDeletedReviewByModuleAndUser(db *sql.DB, moduleId string, userId int) (bool, Review) {
	rows, err := db.Query(`SELECT R.moduleId, R.userId, R.timestamp, R.workload, 
		R.expectedGrade, R.actualGrade, R.difficulty, R.semesterTaken, R.lecturer, 
		R.content, R.suggestion, U.username
		FROM Reviews AS R, Users AS U
		WHERE R.moduleId = $1 AND U.id = $2 AND R.userId = U.id AND R.is_deleted = true`,
		moduleId, userId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var review Review
	isReviewDeleted := false
	for rows.Next() {
		isReviewDeleted = true
		err := rows.Scan(&review.ModuleId, &review.UserId, &review.Timestamp, &review.Workload,
			&review.ExpectedGrade, &review.ActualGrade, &review.Difficulty, &review.SemesterTaken,
			&review.Lecturer, &review.Content, &review.Suggestion, &review.Username)
		if err != nil {
			panic(err)
		}
	}
	if rows.Err() != nil {
		panic(err)
	}
	return isReviewDeleted, review
}

func PostReview(db *sql.DB, moduleId string, userId int, workload int, expectedGrade string,
	actualGrade string, difficulty int, semesterTaken string, lecturer string, content string,
	suggestion string) error {

	exist, _ := GetReviewByModuleAndUser(db, moduleId, userId)
	if exist {
		fmt.Println("Review already exists in the db")
		return errors.New("User has already reviewed this module.")
	}

	fmt.Println("Posting review...")
	fmt.Println(moduleId, userId, content)
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error in initializing db: ", err.Error())
		return errors.New("Unable to begin database transaction")
	}
	defer tx.Rollback()

	isDeleted, _ := GetDeletedReviewByModuleAndUser(db, moduleId, userId)
	if isDeleted {
		_, err = tx.Exec(`DELETE FROM Reviews
		WHERE moduleId = $1 AND userId = $2`,
			moduleId, userId)
		if err != nil {
			panic(err)
		}
	}

	_, err = tx.Exec(`INSERT INTO Reviews (moduleId, userId, workload, expectedGrade, actualGrade, 
		difficulty, semesterTaken, lecturer, content, suggestion) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, moduleId, userId, workload,
		expectedGrade, actualGrade, difficulty, semesterTaken, lecturer, content, suggestion)
	if err != nil {
		fmt.Println("Error in inserting review into db: ", err.Error())
		return errors.New("Review data is malformed.")
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error in commiting review posted: ", err.Error())
		return errors.New("Unable to commit transaction")
	}
	fmt.Println("Posted...")

	return nil
}

func EditReviewByModuleAndUser(db *sql.DB, moduleId string, userId int, workload *int,
	expectedGrade *string, actualGrade *string, difficulty *int, semesterTaken *string,
	lecturer *string, content *string, suggestion *string) error {

	exist, _ := GetReviewByModuleAndUser(db, moduleId, userId)
	if !exist {
		fmt.Println("Review by user does not exist in the db")
		return errors.New("User has never reviewed this module or review has been deleted.")
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.New("Unable to begin database transaction")
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE Reviews 
		SET workload = COALESCE($1, workload), expectedGrade = COALESCE($2, expectedGrade),
		actualGrade = COALESCE($3, actualGrade), difficulty = COALESCE($4, difficulty),
		semesterTaken = COALESCE($5, semesterTaken), lecturer = COALESCE($6, lecturer), 
		content = COALESCE($7, content), suggestion = COALESCE($8, suggestion)
		WHERE moduleId = $9 AND userId = $10 AND is_deleted = false`,
		workload, expectedGrade, actualGrade, difficulty, semesterTaken, lecturer, content,
		suggestion, moduleId, userId)

	if err != nil {
		return errors.New("Some fields have invalid values")
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("Unable to commit transaction")
	}

	return nil
}

func DeleteReviewByModuleAndUser(db *sql.DB, moduleId string, userId int) error {

	exist, _ := GetReviewByModuleAndUser(db, moduleId, userId)
	if !exist {
		fmt.Println("Review by user does not exist in the db")
		return errors.New("User has never reviewed this module or review has been deleted.")
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.New("Unable to begin database transaction")
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE Reviews 
		SET is_deleted = true
		WHERE moduleId = $1 AND userId = $2 AND is_deleted = false`,
		moduleId, userId)

	if err != nil {
		return errors.New("Some fields have invalid values")
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("Unable to commit transaction")
	}

	return nil
}
