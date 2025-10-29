package database

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"sync"

	"github.com/networkgcorefullcode/ssm/factory"
)

// in this package we define functions to generate secure data and store this
// data in the database using best practices.
var SecretStore map[string]string

type Secret struct {
	Secret    string `bson:"encrypted_data"`
	ServiceID string `bson:"service_id"`
}

type dbContext struct {
	GenMutex sync.Mutex
}

var DbContext dbContext = dbContext{
	GenMutex: sync.Mutex{},
}

func GenSecrets() error {
	passwordUDM := generateSecurePassword(16)
	passwordWebconsole := generateSecurePassword(16)

	if SecretStore == nil {
		SecretStore = make(map[string]string)
	}

	// encrypt this secrets before storing in DB

	// secretUDM, err := encryptSecret(passwordUDM)
	// if err != nil {
	// 	return err
	// }
	// secretWebconsole, err := encryptSecret(passwordWebconsole)
	// if err != nil {
	// 	return err
	// }

	secretUDM := Secret{
		Secret:    passwordUDM,
		ServiceID: "udm",
	}
	secretWebconsole := Secret{
		Secret:    passwordWebconsole,
		ServiceID: "webconsole",
	}

	go InsertData(Client, factory.SsmConfig.Configuration.Mongodb.DBName, CollSecret, secretUDM)
	go InsertData(Client, factory.SsmConfig.Configuration.Mongodb.DBName, CollSecret, secretWebconsole)

	return nil
}

// func encryptSecret(secret string) (EncryptedSecret, error) {
// 	password, err := base64.StdEncoding.DecodeString(secret)
// 	if err != nil {
// 		return EncryptedSecret{}, err
// 	}
// 	session := mgrpkcs11.GetSession()
// 	keyHandle, err := pkcs11mgr.FindKeyLabelReturnRandom(constants.LABEL_ENCRYPTION_KEY_INTERNAL_AES256, *session)
// 	if err != nil {
// 		return EncryptedSecret{}, err
// 	}
// 	iv := make([]byte, 16)
// 	if err := safe.RandRead(iv); err != nil {
// 		logger.AppLog.Errorf("Failed to generate IV: %v", err)
// 		return EncryptedSecret{}, err
// 	}
// 	encrypted, err := pkcs11mgr.EncryptKey(keyHandle, iv, password, pkcs11.CKM_AES_CBC_PAD, *session)
// 	if err != nil {
// 		return EncryptedSecret{}, err
// 	}
// 	attrs, err := pkcs11mgr.GetObjectAttributes(keyHandle, *session)
// 	if err != nil {
// 		return EncryptedSecret{}, err
// 	}
// 	mgrpkcs11.LogoutSession(session)
// 	return EncryptedSecret{
// 		EncryptedData: base64.StdEncoding.EncodeToString(encrypted),
// 		IV:            base64.StdEncoding.EncodeToString(iv),
// 		Id:            attrs.Id,
// 		KeyLabel:      constants.LABEL_ENCRYPTION_KEY_INTERNAL_AES256,
// 	}, nil
// }

func generateSecurePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	password := make([]byte, length)
	for i := range password {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[randomIndex.Int64()]
	}

	return base64.StdEncoding.EncodeToString(password)
}
