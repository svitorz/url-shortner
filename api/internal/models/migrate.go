package models

import (
	"fmt"
	"svitorz/url-shortner/internal/repository"
)

func LoadMigrations() {
	rep := repository.DB
	var err error
	if err = rep.Migrator().CreateTable(&User{}); err != nil {
		fmt.Println("Erro ao realizar a migration da tabela user", err)
	}

	if err = rep.Migrator().CreateTable(&Link{}); err != nil {
		fmt.Println("Erro ao realizar a migration da tabela links", err)
	}

	tables, err := rep.Migrator().GetTables()

	for _, name := range tables {
		fmt.Println(name)
	}
}
