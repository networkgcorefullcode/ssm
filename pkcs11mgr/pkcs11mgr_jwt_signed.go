package pkcs11mgr

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
)

var jwtPrivateKey, jwtPublicKey pkcs11.ObjectHandle

// JWTHeader represents the JWT header
type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// JWTPayload represents the JWT payload/claims
type JWTPayload struct {
	Iss string `json:"iss,omitempty"` // Issuer
	Sub string `json:"sub,omitempty"` // Subject
	Aud string `json:"aud,omitempty"` // Audience
	Exp int64  `json:"exp,omitempty"` // Expiration time
	Nbf int64  `json:"nbf,omitempty"` // Not before
	Iat int64  `json:"iat,omitempty"` // Issued at
	Jti string `json:"jti,omitempty"` // JWT ID
}

// GetJWTPrivateKey returns the JWT private key handle
func GetJWTPrivateKey() pkcs11.ObjectHandle {
	return jwtPrivateKey
}

// GetJWTPublicKey returns the JWT public key handle
func GetJWTPublicKey() pkcs11.ObjectHandle {
	return jwtPublicKey
}

// InitJWTKey initializes the JWT signing key by finding it in the HSM
func InitJWTKey(s *Session) error {
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, constants.JWTKeyLabel),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
	}

	if err := s.Ctx.FindObjectsInit(s.Handle, template); err != nil {
		logger.AppLog.Errorf("Failed to initialize JWT key search: %v", err)
		return err
	}
	defer s.Ctx.FindObjectsFinal(s.Handle)

	obj, _, err := s.Ctx.FindObjects(s.Handle, 1)
	if err != nil {
		logger.AppLog.Errorf("Failed to find JWT private key: %v", err)
		return err
	}

	if len(obj) == 0 {
		logger.AppLog.Warn("JWT private key not found, will generate new key pair")
		return generateJWTKeyPair(s)
	}

	jwtPrivateKey = obj[0]

	// Load the corresponding public key
	pubTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, constants.JWTKeyLabel),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
	}

	if err := s.Ctx.FindObjectsInit(s.Handle, pubTemplate); err != nil {
		logger.AppLog.Errorf("Failed to initialize JWT public key search: %v", err)
		return err
	}
	defer s.Ctx.FindObjectsFinal(s.Handle)

	pubObj, _, err := s.Ctx.FindObjects(s.Handle, 1)
	if err != nil {
		logger.AppLog.Errorf("Failed to find JWT public key: %v", err)
		return err
	}

	if len(pubObj) > 0 {
		jwtPublicKey = pubObj[0]
		logger.AppLog.Info("JWT key pair loaded successfully")
	} else {
		logger.AppLog.Warn("JWT public key not found")
	}

	return nil
}

// generateJWTKeyPair generates a new RSA key pair for JWT signing
func generateJWTKeyPair(s *Session) error {
	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, constants.JWTKeyLabel),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS_BITS, 2048),
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, []byte{1, 0, 1}),
	}

	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, constants.JWTKeyLabel),
	}

	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_KEY_PAIR_GEN, nil)}

	pubKey, privKey, err := s.Ctx.GenerateKeyPair(
		s.Handle,
		mechanism,
		publicKeyTemplate,
		privateKeyTemplate,
	)

	if err != nil {
		logger.AppLog.Errorf("Failed to generate JWT key pair: %v", err)
		return err
	}

	jwtPrivateKey = privKey
	jwtPublicKey = pubKey
	logger.AppLog.Infof("Generated new JWT key pair - Public: %d, Private: %d", pubKey, privKey)
	return nil
}

// SignJWT creates and signs a JWT token using the HSM private key
func SignJWT(s *Session, payload JWTPayload) (string, error) {
	// Set default values if not provided
	if payload.Iat == 0 {
		payload.Iat = time.Now().Unix()
	}
	if payload.Exp == 0 {
		payload.Exp = time.Now().Add(time.Hour * 24).Unix() // 24 hours default
	}

	// Create JWT header
	header := JWTHeader{
		Alg: "RS256",
		Typ: "JWT",
	}

	// Encode header
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %v", err)
	}
	headerEncoded := hex.EncodeToString(headerJSON)

	// Encode payload
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}
	payloadEncoded := hex.EncodeToString(payloadJSON)

	// Create signing input
	signingInput := headerEncoded + "." + payloadEncoded

	// Hash the signing input
	hasher := sha256.New()
	hasher.Write([]byte(signingInput))
	hashed := hasher.Sum(nil)

	// Sign with HSM
	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA256_RSA_PKCS, nil)}

	err = s.Ctx.SignInit(s.Handle, mechanism, jwtPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to init signing: %v", err)
	}

	signature, err := s.Ctx.Sign(s.Handle, hashed)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %v", err)
	}

	// Encode signature
	signatureEncoded := hex.EncodeToString(signature)

	// Return complete JWT
	jwt := signingInput + "." + signatureEncoded
	logger.AppLog.Infof("JWT signed successfully, length: %d", len(jwt))

	return jwt, nil
}

// VerifyJWT verifies a JWT token using the HSM public key
func VerifyJWT(s *Session, token string) (*JWTPayload, error) {
	// Split the token
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	headerEncoded, payloadEncoded, signatureEncoded := parts[0], parts[1], parts[2]

	// Decode signature
	signature, err := hex.DecodeString(signatureEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Create signing input for verification
	signingInput := headerEncoded + "." + payloadEncoded

	// Hash the signing input
	hasher := sha256.New()
	hasher.Write([]byte(signingInput))
	hashed := hasher.Sum(nil)

	// Verify signature with HSM
	mechanism := []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA256_RSA_PKCS, nil)}

	err = s.Ctx.VerifyInit(s.Handle, mechanism, jwtPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to init verification: %v", err)
	}

	err = s.Ctx.Verify(s.Handle, hashed, signature)
	if err != nil {
		return nil, fmt.Errorf("JWT signature verification failed: %v", err)
	}

	// Decode and parse payload
	payloadJSON, err := hex.DecodeString(payloadEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	var payload JWTPayload
	err = json.Unmarshal(payloadJSON, &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payload: %v", err)
	}

	// Check expiration
	if payload.Exp > 0 && time.Now().Unix() > payload.Exp {
		return nil, fmt.Errorf("JWT token has expired")
	}

	// Check not before
	if payload.Nbf > 0 && time.Now().Unix() < payload.Nbf {
		return nil, fmt.Errorf("JWT token is not yet valid")
	}

	logger.AppLog.Info("JWT verified successfully")
	return &payload, nil
}

// CreateStandardJWT creates a JWT with standard claims
func CreateStandardJWT(s *Session, issuer, subject, audience string, expirationHours int) (string, error) {
	now := time.Now()
	payload := JWTPayload{
		Iss: issuer,
		Sub: subject,
		Aud: audience,
		Iat: now.Unix(),
		Exp: now.Add(time.Duration(expirationHours) * time.Hour).Unix(),
		Jti: fmt.Sprintf("%d", now.UnixNano()), // Simple unique ID
	}

	return SignJWT(s, payload)
}
