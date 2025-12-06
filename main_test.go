package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")

	if err != nil {
		t.Fatal("error db connection:", err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Fatal("db ping error:", err)
	}

	clientID := 1

	res, err := selectClient(db, clientID)

	if err != nil {
		t.Fatal("selectClient вернула ошибку, хотя должна была работать:", err)
	}

	assert.Equal(t, clientID, res.ID, "ID должен совпадать с переданным clientID")
	assert.NotEmpty(t, res.FIO, "FIO не должен быть пустым")
	assert.NotEmpty(t, res.Login, "Login не должен быть пустым")
	assert.NotEmpty(t, res.Birthday, "Birthday не должен быть пустым")
	assert.NotEmpty(t, res.Email, "Email не должен быть пустым")

}

func Test_SelectClient_WhenNoClient(t *testing.T) {

	db, err := sql.Open("sqlite", "demo.db")

	if err != nil {
		t.Fatal("error db connection:", err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Fatal("db ping error:", err)
	}

	clientID := -1

	res, err := selectClient(db, clientID)

	require.Error(t, err)
	require.ErrorIs(t, sql.ErrNoRows, err)

	assert.Zero(t, res.ID, "ID должен быть пустым")
	assert.Empty(t, res.FIO, "FIO должен быть пустым")
	assert.Empty(t, res.Login, "Login должен быть пустым")
	assert.Empty(t, res.Birthday, "Birthday должен быть пустым")
	assert.Empty(t, res.Email, "Email должен быть пустым")

}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	db, err := sql.Open("sqlite", "demo.db")

	if err != nil {
		t.Fatal("error db connection:", err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Fatal("db ping error:", err)
	}

	id, err := insertClient(db, cl)

	cl.ID = id

	require.NotZero(t, id)
	require.NoError(t, err)

	selectCl, err := selectClient(db, id)

	require.NoError(t, err)
	require.Equal(t, cl.ID, selectCl.ID)
	require.Equal(t, cl.FIO, selectCl.FIO)
	require.Equal(t, cl.Login, selectCl.Login)
	require.Equal(t, cl.Birthday, selectCl.Birthday)
	require.Equal(t, cl.Email, selectCl.Email)

}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")

	if err != nil {
		t.Fatal("error db connection:", err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Fatal("db ping error:", err)
	}

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	id, err := insertClient(db, cl)

	require.NotZero(t, id)
	require.NoError(t, err)

	selectCl, err := selectClient(db, id)
	require.NoError(t, err)

	err = deleteClient(db, selectCl.ID)
	require.NoError(t, err)

	_, err = selectClient(db, selectCl.ID)

	assert.ErrorIs(t, sql.ErrNoRows, err)
}
