// +build integration

package mysql

import (
	sharesecret "github.com/bernardosecades/sharesecret/internal"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var mr sharesecret.SecretRepository

func init() {
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	mr = NewMySQLSecretRepository(dbName, dbUser, dbPass, dbHost, dbPort)
}

func TestMySQLSecretRepository_All_Methods(t *testing.T) {

	//layout := "2006-01-02 15:04:05"
	//str := "2046-01-02 15:04:05"
	//tm,_ := time.Parse(layout, str)

	tm := time.Now().UTC().Add(-1 * time.Hour)

	r1, err := mr.CreateSecret("this is a test", true, tm)
	assert.Nil(t, err)
	assert.NotNil(t, r1)

	/*
		r2, err := mr.GetSecret(r1.ID)
		assert.Nil(t, err)
		assert.NotNil(t, r2)

		customPass, err := mr.HasSecretWithCustomPwd(r2.ID)
		assert.Nil(t, err)
		assert.True(t, customPass)

		err := mr.RemoveSecret(r1.ID)
		assert.Nil(err)
	*/
}
