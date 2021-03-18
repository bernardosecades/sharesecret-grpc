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

func TestMySQLSecretRepositoryCreateAndReadSecretNoExpired(t *testing.T) {

	tm := time.Now().UTC().Add(time.Hour)
	r1, err1 := mr.CreateSecret("this is a test create and read secret not expired", true, tm)

	assert.Nil(t, err1)
	assert.NotNil(t, r1)

	r2, err2 := mr.HasSecretWithCustomPwd(r1.ID)

	assert.Nil(t, err2)
	assert.True(t, r2)

	r3, err3 := mr.GetSecret(r1.ID)

	assert.Nil(t, err3)
	assert.Equal(t, "this is a test create and read secret not expired", r3.Content)

	err4 := mr.RemoveSecret(r3.ID)

	assert.Nil(t, err4)
}

func TestMySQLSecretRepositoryCreateAndReadSecretExpired(t *testing.T) {

	tm := time.Now().UTC().Add(-1 * time.Hour)
	r1, err1 := mr.CreateSecret("this is a test create and read secret not expired", true, tm)

	assert.Nil(t, err1)
	assert.NotNil(t, r1)

	r2, err2 := mr.HasSecretWithCustomPwd(r1.ID)

	assert.NotNil(t, err2)
	assert.False(t, r2)

	_, err3 := mr.GetSecret(r1.ID)

	assert.NotNil(t, err3)

	r4, err4 := mr.RemoveSecretsExpired()

	assert.Nil(t, err4)
	assert.Equal(t, int64(1), r4)
}
