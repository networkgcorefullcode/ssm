package database

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

// in this package we define functions to generate secure data and store this
// data in the database using best practices.
var secretStore map[string]string

type EncryptedSecret struct {
	EncryptedData string `bson:"encrypted_data"`
	IV            string `bson:"iv"`
	Id            int32  `bson:"id"`
	KeyLabel      string `bson:"key_label"`
}

func SetSecretMap(secrets map[string]string) {
	secretStore = secrets
}

func GenSecrets() error {
	passwordUDM := generateSecurePassword(16)
	passwordWebconsole := generateSecurePassword(16)

	if secretStore == nil {
		secretStore = make(map[string]string)
	}

	// encrypt this secrets before storing in DB

	secretUDM, err := encryptSecret(passwordUDM)
	if err != nil {
		return err
	}
	secretWebconsole, err := encryptSecret(passwordWebconsole)
	if err != nil {
		return err
	}

	InsertData(Client, factory.SsmConfig.Configuration.Mongodb.DBName, CollSecret, secretUDM)
	InsertData(Client, factory.SsmConfig.Configuration.Mongodb.DBName, CollSecret, secretWebconsole)

	return nil
}

func encryptSecret(secret string) (EncryptedSecret, error) {
	password, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return EncryptedSecret{}, err
	}
	session := mgrpkcs11.GetSession()
	keyHandle, err := pkcs11mgr.FindKeyLabelReturnRandom(constants.LABEL_ENCRYPTION_KEY_INTERNAL_AES256, *session)
	if err != nil {
		return EncryptedSecret{}, err
	}
	iv := make([]byte, 16)
	if err := safe.RandRead(iv); err != nil {
		logger.AppLog.Errorf("Failed to generate IV: %v", err)
		return EncryptedSecret{}, err
	}
	encrypted, err := pkcs11mgr.EncryptKey(keyHandle, iv, password, pkcs11.CKM_AES_CBC_PAD, *session)
	if err != nil {
		return EncryptedSecret{}, err
	}
	attrs, err := pkcs11mgr.GetObjectAttributes(keyHandle, *session)
	if err != nil {
		return EncryptedSecret{}, err
	}
	mgrpkcs11.LogoutSession(session)
	return EncryptedSecret{
		EncryptedData: base64.StdEncoding.EncodeToString(encrypted),
		IV:            base64.StdEncoding.EncodeToString(iv),
		Id:            attrs.Id,
	}, nil
}

func generateSecurePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	password := make([]byte, length)
	for i := range password {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[randomIndex.Int64()]
	}

	return base64.StdEncoding.EncodeToString(password)
}
