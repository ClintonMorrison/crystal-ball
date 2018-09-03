package main

import (
  "database/sql"
  "fmt"
  _ "github.com/lib/pq"
)

func GetHandle() *sql.DB {
  connStr := fmt.Sprintf(
    "postgres://%s:%s@%s/%s",
     DB_USER,
     DB_PASSWORD,
     DB_HOST,
     DB_NAME)

  fmt.Println("OPENING!")
  db, err := sql.Open("postgres", connStr)

  if err != nil {
    panic("could not connect")
  }

  return db
}



