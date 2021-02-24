package usr

import (
	"github.com/mcaci/login/db"
)

const users = "users"

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Save(d *db.Database, user *User) error {
	pipe := d.Client.TxPipeline()
	pipe.HSet(db.Ctx, user.Username, "password", user.Password)
	_, err := pipe.Exec(db.Ctx)
	if err != nil {
		return err
	}
	return nil
}

func Get(d *db.Database, username string) (*User, error) {
	pipe := d.Client.TxPipeline()
	password := pipe.HGet(db.Ctx, username, "password")
	_, err := pipe.Exec(db.Ctx)
	if err != nil {
		return nil, err
	}
	if password == nil {
		return nil, db.ErrNil
	}
	return &User{
		Username: username,
		Password: password.Val(),
	}, nil
}
