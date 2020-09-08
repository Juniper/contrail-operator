package client

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Postgres struct {
	db *pg.DB
}

func New(address, user, password, dbName string) (*Postgres, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s", user, password, address, dbName)
	opt, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)
	return &Postgres{
		db: db,
	}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) ping(ctx context.Context) error {
	return p.db.Ping(ctx)
}

type TestUser struct {
	ID   int64
	Name string
}

func (p *Postgres) CreateTestTable(ctx context.Context) error {
	if err := p.ping(ctx); err != nil {
		return err
	}
	model := &TestUser{}
	return p.db.CreateTable(model, &orm.CreateTableOptions{
		Temp:        false,
		IfNotExists: true,
	})
}

func (p *Postgres) InsertTestUser(ctx context.Context, id int64, testData string) error {
	if err := p.ping(ctx); err != nil {
		return err
	}
	user := &TestUser{
		ID:   id,
		Name: testData,
	}
	return p.db.Insert(user)
}

func (p *Postgres) GetTestUserName(ctx context.Context, id int64) (string, error) {
	if err := p.ping(ctx); err != nil {
		return "", err
	}
	user := &TestUser{
		ID: id,
	}
	err := p.db.Select(user)
	if err != nil {
		return "", err
	}
	return user.Name, nil
}
