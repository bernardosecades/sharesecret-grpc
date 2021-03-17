package sharesecret

import (
	"encoding/hex"
	"errors"
	"github.com/bernardosecades/sharesecret/internal/util"
	"time"
)

// All errors reported by the service
var (
	ErrSecretNotFound = errors.New("it either never existed or has already been viewed")
	ErrNoPassRequired = errors.New("the password is not required")
	ErrMissingPass    = errors.New("you need a password to see the secret")
	ErrEmptyContent   = errors.New("empty content")
	ErrTextTooLong    = errors.New("text too long")
	ErrPassTooLong    = errors.New("password too long")
	ErrPassToDecrypt  = errors.New("error password to decrypt")
	ErrToEncrypt      = errors.New("error to encrypt")
)

type SecretService interface {
	GetContentSecret(id string, password string) (string, error)
	CreateSecret(rawContent string, password string) (Secret, error)
}

type secretService struct {
	repository SecretRepository
	key        string
	defaultPwd string
}

func NewSecretService(r SecretRepository, key string, defaultPwd string) SecretService {

	if len(key) != 32 {
		panic("key secret should have 32 bytes")
	}

	return &secretService{r, key, defaultPwd}
}

func (s *secretService) GetContentSecret(id string, password string) (string, error) {

	hasPass, err := s.hasSecretWithCustomPwd(id)

	if err != nil {
		return "", ErrSecretNotFound
	}

	if hasPass && len(password) == 0 {
		return "", ErrMissingPass
	}

	if !hasPass && len(password) > 0 {
		return "", ErrNoPassRequired
	}

	if len(password) == 0 {
		password = s.defaultPwd
	}

	secret, err := s.getSecret(id)
	if err != nil {
		return "", ErrSecretNotFound
	}

	content, err := s.decryptContentSecret(secret.Content, password)

	if err != nil {
		return "", ErrPassToDecrypt
	}

	return content, nil
}

func (s *secretService) CreateSecret(rawContent string, password string) (Secret, error) {

	if len(rawContent) == 0 {
		return Secret{}, ErrEmptyContent
	}

	if len(rawContent) > 10000 {
		return Secret{}, ErrTextTooLong
	}

	if len(password) > 32 {
		return Secret{}, ErrPassTooLong
	}

	customPwd := true
	if len(password) == 0 {
		customPwd = false
		password = s.defaultPwd
	}

	content, err := s.encryptContentSecret(rawContent, password)

	if err != nil {
		return Secret{}, ErrToEncrypt
	}

	expire := time.Now().UTC().AddDate(0, 0, 5)
	secret, err := s.repository.CreateSecret(content, customPwd, expire)
	if err != nil {
		return Secret{}, err
	}

	return secret, nil
}

func (s secretService) hasSecretWithCustomPwd(id string) (bool, error) {

	return s.repository.HasSecretWithCustomPwd(id)
}

func (s secretService) getSecret(id string) (Secret, error) {
	secret, err := s.repository.GetSecret(id)
	if err != nil {
		return Secret{}, err
	}

	err = s.repository.RemoveSecret(id)

	if err != nil {
		return Secret{}, err
	}

	return secret, nil
}

func (s *secretService) decryptContentSecret(content string, password string) (string, error) {
	decodeContent, _ := hex.DecodeString(content)
	key := []byte(s.key)
	copy(key[:], password)
	decryptContent, err := util.Decrypt(key, decodeContent)

	if err != nil {
		return "", err
	}

	return string(decryptContent), nil
}

func (s *secretService) encryptContentSecret(content string, password string) (string, error) {
	key := []byte(s.key)
	copy(key[:], password)
	encryptContent, err := util.Encrypt(key, []byte(content))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encryptContent), nil
}
