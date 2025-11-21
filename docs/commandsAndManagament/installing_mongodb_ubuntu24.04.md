# MongoDB Installation and ReplicaSet Configuration on Ubuntu 24.04

This guide provides step-by-step commands to install MongoDB Community Edition on Ubuntu 24.04 and configure it as a ReplicaSet.

## System Information
- **OS**: Ubuntu 24.04 LTS
- **MongoDB Version**: 7.0 (Community Edition)
- **IP Address**: 192.168.12.16
- **ReplicaSet Name**: rs0

---

## 1. Install MongoDB

### Import MongoDB GPG Key
```bash
curl -fsSL https://www.mongodb.org/static/pgp/server-7.0.asc | \
   sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg --dearmor
```

### Add MongoDB Repository
```bash
echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | \
   sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list
```

### Update Package Database
```bash
sudo apt-get update
```

### Install MongoDB Packages
```bash
sudo apt-get install -y mongodb-org
```

### Pin MongoDB Version (Optional - Prevents Automatic Updates)
```bash
echo "mongodb-org hold" | sudo dpkg --set-selections
echo "mongodb-org-database hold" | sudo dpkg --set-selections
echo "mongodb-org-server hold" | sudo dpkg --set-selections
echo "mongodb-mongosh hold" | sudo dpkg --set-selections
echo "mongodb-org-mongos hold" | sudo dpkg --set-selections
echo "mongodb-org-tools hold" | sudo dpkg --set-selections
```

---

## 2. Configure MongoDB for ReplicaSet

### Edit MongoDB Configuration File
```bash
sudo nano /etc/mongod.conf
```

### Update the Configuration
Replace or modify the following sections:

```yaml
# network interfaces
net:
  port: 27017
  bindIp: 192.168.12.16,127.0.0.1  # Listen on specific IP and localhost

# replication
replication:
  replSetName: "rs0"
```

**Complete Configuration Example:**
```yaml
# mongod.conf

# for documentation of all options, see:
#   http://docs.mongodb.org/manual/reference/configuration-options/

# Where and how to store data.
storage:
  dbPath: /var/lib/mongodb
  journal:
    enabled: true

# where to write logging data.
systemLog:
  destination: file
  logAppend: true
  path: /var/log/mongodb/mongod.log

# network interfaces
net:
  port: 27017
  bindIp: 192.168.12.16,127.0.0.1

# how the process runs
processManagement:
  timeZoneInfo: /usr/share/zoneinfo

# security
security:
  authorization: disabled  # Enable after initial setup

# replication
replication:
  replSetName: "rs0"
```

---

## 3. Start MongoDB Service

### Enable MongoDB to Start on Boot
```bash
sudo systemctl enable mongod
```

### Start MongoDB Service
```bash
sudo systemctl start mongod
```

### Check MongoDB Status
```bash
sudo systemctl status mongod
```

### View MongoDB Logs (if needed)
```bash
sudo tail -f /var/log/mongodb/mongod.log
```

---

## 4. Initialize ReplicaSet

### Connect to MongoDB Shell
```bash
mongosh --host 192.168.12.16
```

### Initialize the ReplicaSet
Run this command in the MongoDB shell:

```javascript
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "192.168.12.16:27017" }
  ]
})
```

```javascript
// Ver configuración actual
rs.conf()

// Reconfigurar con el host correcto
cfg = rs.conf()
cfg.members[0].host = "192.168.12.16:27017"
rs.reconfig(cfg)
```

### Verify ReplicaSet Status
```javascript
rs.status()
```

### Check ReplicaSet Configuration
```javascript
rs.conf()
```

---

## 5. Verify Installation

### Check ReplicaSet Member State
```javascript
rs.isMaster()
```

You should see:
```json
{
  "ismaster": true,
  "topologyVersion": { ... },
  "hosts": [ "192.168.12.16:27017" ],
  "setName": "rs0",
  "primary": "192.168.12.16:27017",
  ...
}
```

### Create Test Database
```javascript
use testdb
db.testCollection.insertOne({ test: "Hello ReplicaSet" })
db.testCollection.find()
```

### Exit MongoDB Shell
```javascript
exit
```

---

## 6. Firewall Configuration (if UFW is enabled)

```bash
sudo ufw allow from 192.168.12.0/24 to any port 27017
sudo ufw reload
```

---

## 7. Connection String for Applications

Use this connection string in your applications:

```
mongodb://192.168.12.16:27017/?replicaSet=rs0
```

For the SSM service (`ssmConfig.yml`):
```yaml
mongodb:
  name: "ssm_db"
  url: "mongodb://192.168.12.16:27017/?replicaSet=rs0"
  dbName: "secure_storage"
```

---

## 8. Useful Commands

### Restart MongoDB
```bash
sudo systemctl restart mongod
```

### Stop MongoDB
```bash
sudo systemctl stop mongod
```

### Check MongoDB Logs
```bash
sudo journalctl -u mongod -f
```

### MongoDB Shell Connection
```bash
mongosh --host 192.168.12.16
```

### Check ReplicaSet Status (from shell)
```bash
mongosh --host 192.168.12.16 --eval "rs.status()"
```

---

## 9. Security Recommendations (Post-Setup)

### Create Admin User
```javascript
use admin
db.createUser({
  user: "admin",
  pwd: "SecurePassword123!",
  roles: [ { role: "userAdminAnyDatabase", db: "admin" }, "readWriteAnyDatabase" ]
})
```

### Enable Authentication
Edit `/etc/mongod.conf`:
```yaml
security:
  authorization: enabled
```

Restart MongoDB:
```bash
sudo systemctl restart mongod
```

### Connect with Authentication
```bash
mongosh --host 192.168.12.16 -u admin -p SecurePassword123! --authenticationDatabase admin
```

---

## Troubleshooting

### Check if MongoDB is Running
```bash
sudo systemctl status mongod
ps aux | grep mongod
```

### Check Port Binding
```bash
sudo netstat -tulpn | grep 27017
```

### Check MongoDB Logs for Errors
```bash
sudo cat /var/log/mongodb/mongod.log | grep ERROR
```

### Test Network Connectivity
```bash
telnet 192.168.12.16 27017
```

### ReplicaSet Not Initializing
- Ensure `bindIp` includes the IP address (192.168.12.16)
- Check firewall rules
- Verify DNS/hostname resolution

---

## Summary

You now have:
- ✅ MongoDB 7.0 installed on Ubuntu 24.04
- ✅ Configured to listen on 192.168.12.16
- ✅ ReplicaSet named `rs0` initialized
- ✅ Ready for production use with the connection string: `mongodb://192.168.12.16:27017/?replicaSet=rs0`
