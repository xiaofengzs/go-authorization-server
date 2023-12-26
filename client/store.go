package client

import (
	"github.com/ory/fosite"
	"github.com/xiaofengzs/go-authorization-server/database"
)

type Store struct {
	*database.MysqlProvider
}

func NewClientStore(mysql *database.MysqlProvider) *Store {
	return &Store{mysql}
}

func (store *Store) Save(client fosite.DefaultClient) error {
	sql := "insert into client('', '', '', '', '') values ($1, $2, $3, $4, $5)"
	result, err := store.Exec(sql)
	if err != nil {
		return err
	}

	result.LastInsertId()
	return nil
}

func (store *Store) Get(clientId string) (fosite.Client, error) {
	sql := "select * from client where id = $1"
	rows, err := store.Query(sql, clientId)
	if err != nil {
		return nil, err
	}
	client := fosite.DefaultClient{}
	if err := rows.Scan(&client.ID, client.Secret, client.GrantTypes, client.Scopes); err != nil {
		return nil, err
	}

	return &client, nil
}
