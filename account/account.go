package account

type Account struct {
	ID int64
	name string
}

type Credentials struct {
	ID int64 `json:"id"`
	Key string `json:"key"`
}

func (account *Account) SetName(Name string) (*Account) {

	// do validation
	// enter into database

	account.name = Name
	return account
}

func (account *Account) Name() string {
	return account.name
}

func validateName(name string) error {
	return nil
}

func NewAccount(name string) (acc *Account, credentials *Credentials, err error) {

	err = validateName(name)
	if err != nil {
		return
	}


	// do validation on name

	acc = &Account {
		int64(32424),
		"Billy",
	}

	credentials = &Credentials {
		ID:acc.ID,
		Key:"key",
	}

	return
}

func RetrieveAccount(accountID int64, key string) (*Account, error) {
	acc := &Account {
		int64(32424),
		"Billy",
	}

	return acc, nil
}