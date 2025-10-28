package database

import (
	"context"
	"time"

	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = &mongo.Client{}
var mgrpkcs11 *pkcs11mgr.Manager

func initDatabase(client_pass *mongo.Client, url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	var err error
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return err
	}
	*client_pass = *client
	return nil
}

func init() {
	initDatabase(Client, factory.SsmConfig.Configuration.Mongodb.Url)
}

func SetPKCS11Manager(mgr *pkcs11mgr.Manager) {
	mgrpkcs11 = mgr
}
