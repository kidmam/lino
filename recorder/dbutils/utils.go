package dbutils

import (
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	"github.com/lino-network/lino/recorder/errors"
)

const (
	// CoinStrLength is the length of the coin string
	CoinStrLength = 64
)

// ExecAffectingOneRow executes a given statement, expecting one row to be affected.
func ExecAffectingOneRow(stmt *sql.Stmt, args ...interface{}) (sql.Result, errors.Error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return r, errors.Internalf("ExecAffectingOneRow: failed to execute statement [%v]", stmt).TraceCause(err, "")
	}
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return r, errors.Internalf("ExecAffectingOneRow: can't get rows affected for statement [%v]", stmt).TraceCause(err, "")
	} else if rowsAffected != 1 {
		return r, errors.Internalf("ExecAffectingOneRow: expect 1, but got [%d] row affected for statement [%v]", rowsAffected, stmt)
	}
	return r, nil
}

// Exec executes a given statement
func Exec(stmt *sql.Stmt, args ...interface{}) (sql.Result, errors.Error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return r, errors.Internalf("Exec: failed to execute statement [%v]", stmt).TraceCause(err, "")
	}
	return r, nil
}

// RowScanner is implemented by sql.Row and sql.Rows
type RowScanner interface {
	Scan(dest ...interface{}) error
}

// PrepareStmts will attempt to prepare each unprepared
// query on the database. If one fails, the function returns
// with an error.
func PrepareStmts(service string, db *sql.DB, unprepared map[string]string) (map[string]*sql.Stmt, errors.Error) {
	prepared := map[string]*sql.Stmt{}
	for k, v := range unprepared {
		stmt, err := db.Prepare(v)
		if err != nil {
			return nil, errors.UnablePrepareStatement(fmt.Sprintf("service: %s can't prepare %v statement", service, stmt)).TraceCause(err, "")
		}
		prepared[k] = stmt
	}

	return prepared, nil
}

// PadNumberStrWithZero pad a number string with zero
func PadNumberStrWithZero(number string) (string, errors.Error) {
	if len(number) == 0 {
		return "", errors.NewError(errors.CodeUnablePrepareStatement, "util.PadNumberStrWithZero: found number with zero")
	}
	if len(number) > CoinStrLength {
		return "", errors.NewErrorf(errors.CodeUnablePrepareStatement, "util.PadNumberStrWithZero: cannot pad number larger than %d length", CoinStrLength)
	}
	paddedNumber := fmt.Sprintf("%064s", number)
	return paddedNumber, nil
}

func checkNumberStrIsValid(number string) bool {
	for _, c := range number {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func TrimPaddedZeroFromNumber(number string) string {
	trimmedNumber := strings.TrimLeft(number, "0")
	if len(trimmedNumber) == 0 {
		return "0"
	}
	return trimmedNumber
}