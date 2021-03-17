package sharesecret

import "time"

type SecretRepository interface {
	GetSecret(id string) (Secret, error)
	CreateSecret(content string, customPwd bool, expire time.Time) (Secret, error)
	RemoveSecret(id string) error
	RemoveSecretsExpired() (int64, error)
	HasSecretWithCustomPwd(id string) (bool, error)
}
