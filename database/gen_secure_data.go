package database

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
	"go.mongodb.org/mongo-driver/bson"
)

type UserSecret struct {
	PasswordSecret EncryptedSecret `bson:"encrypted_data"`
	ServiceID      string          `bson:"service_id"`
}

type EncryptedSecret struct {
	EncryptedData       string `bson:"encrypted_data"`
	IV                  string `bson:"iv"`
	Id                  int32  `bson:"id"`
	KeyLabel            string `bson:"key_label"`
	EncryptionAlgorithm uint   `bson:"encryption_algorithm"`
}

type dbContext struct {
	GenMutex sync.Mutex
}

var DbContext dbContext = dbContext{
	GenMutex: sync.Mutex{},
}

func GenSecrets() error {
	if err := genSecret("udm"); err != nil {
		return err
	}
	if err := genSecret("webconsole"); err != nil {
		return err
	}

	return nil
}

func genSecret(serviceID string) error {

	filter := bson.M{"service_id": serviceID}
	_, err := FindOneData(Client, factory.SsmConfig.Configuration.Mongodb.DBName, CollSecret, filter)
	if err == nil {
		logger.AppLog.Infof("User found: %v", err)
		return nil
	}

	password := generateSecurePassword(16)

	// encrypt this secrets before storing in DB
	secret, err := encryptSecret(password)
	if err != nil {
		return err
	}

	userSecret := UserSecret{
		PasswordSecret: secret,
		ServiceID:      serviceID,
	}

	// Create the directory if it doesn't exist
	secretDir := "/tmp/user-secret-ssm"
	if err := os.MkdirAll(secretDir, 0755); err != nil {
		logger.AppLog.Errorf("Failed to create directory %s: %v", secretDir, err)
		return err
	}

	file, err := os.OpenFile("/tmp/user-secret-ssm/user_secrets.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.AppLog.Errorf("Failed to create user_secrets.txt file: %v", err)
		return err
	}
	defer file.Close()

	envContent := fmt.Sprintf(`user_%s = "%s"
	password_secret_%s = "%s"

	`, serviceID, serviceID, serviceID, password)

	_, err = file.WriteString(envContent)
	if err != nil {
		logger.AppLog.Errorf("Failed to write to user_secrets.txt file: %v", err)
	}

	if _, err = InsertData(Client, factory.SsmConfig.Configuration.Mongodb.DBName, CollSecret, userSecret); err != nil {
		return err
	}

	return nil
}

func encryptSecret(secret string) (EncryptedSecret, error) {
	password, err := hex.DecodeString(secret)
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
		EncryptedData: hex.EncodeToString(encrypted),
		IV:            hex.EncodeToString(iv),
		Id:            attrs.Id,
		KeyLabel:      constants.LABEL_ENCRYPTION_KEY_INTERNAL_AES256,
	}, nil
}

func generateSecurePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	password := make([]byte, length)
	for i := range password {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[randomIndex.Int64()]
	}

	return hex.EncodeToString(password)
}
