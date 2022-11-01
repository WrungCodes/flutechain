package database

type Account string

type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

func (t Tx) isReward() bool {
	return t.Data == "reward"
}

func NewTx(from, to Account, value uint, data string) Tx {
	return Tx{Account(from), Account(to), value, data}
}
