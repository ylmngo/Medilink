package data

import (
	"database/sql"
)

type Model struct {
	Um *UserModel
	Fm *FileModel
}

func NewModel(db *sql.DB) *Model {
	model := &Model{
		Um: &UserModel{DB: db},
		Fm: &FileModel{DB: db},
	}

	// user := newUser("Aamer Aijaz", "Hello", "aamerasim45@gmail.com")
	// fmt.Println(user.Email)
	// err := model.um.Insert(user)
	// if err != nil {
	// 	fmt.Printf("Could not insert to Database: %v\n", err)
	// }

	return model
}
