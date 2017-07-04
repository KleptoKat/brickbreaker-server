package account

import (
	"errors"
	"github.com/KleptoKat/brickbreaker-server/database"
	"github.com/chrislonng/starx/log"
	"math/rand"
)

type Account struct {
	ID   int64 `json:"id"`
	Name string `json:"name"`
}

type Credentials struct {
	ID int64 `json:"id,omitempty"`
	Key string `json:"key"`
}

func (account *Account) SetName(Name string) (*Account) {

	// do validation
	// enter into database

	account.Name = Name
	return account
}

func validateName(name string) bool {

	return len(name) >= 1 && len(name) <= 12
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
	}
	return string(b)
}

func NewAccount(name string) (acc *Account, credentials *Credentials, err error) {

	if !validateName(name) {
		return nil, nil, errors.New("Invalid Name.")
	}
	key := RandStringBytesRmndr(128)

	stmt, err := database.DB().Prepare("INSERT INTO Account (Name, AuthKey) VALUES (?, ?);")

	if err != nil {
		log.Error(err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(name, key)
	if err != nil {
		log.Error(err)
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Error(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Error(err)
	}

	if rowCnt == 0 {
		log.Error(err)
	}


	acc = &Account {
		int64(id),
		name,
	}

	credentials = &Credentials {
		ID:acc.ID,
		Key:key,
	}

	return
}


func RetrieveAccount(id int64) (acc *Account) {
	acc = &Account {}
	err := database.DB().QueryRow("select ID, Name from Account where ID = ?", id).Scan(&acc.ID, &acc.Name)

	if err != nil {
		log.Error(err)
		return nil
	}

	return
}

func AuthenticateAccount(id int64, key string) (*Account) {

	acc := &Account {}

	err := database.DB().QueryRow(
		"select ID, Name from Account where ID = ? AND AuthKey = ?",
		id, key).Scan(
		&acc.ID, &acc.Name)

	if err != nil {
		return nil
	}

	return acc
}


func GetNameByID(id int64) (name string) {
	name = ""
	err := database.DB().QueryRow("select Name from Account where ID = ?", id).Scan(&name)

	if err != nil {
		log.Error(err)
		return
	}

	return


}