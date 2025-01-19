package internal

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
	"github.com/josephakayesi/go-cerbos-abac/infra/config"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

var c = config.GetConfig()

type Paseto struct {
	v2 *paseto.V2
}

type ContextWithValueKey string

// GetBrowserFingerPrintKey: returns the key for setting a browser fingerprint in a context
func GetBrowserFingerPrintKey() ContextWithValueKey {
	return ContextWithValueKey("BrowserFingerprint")
}

func GetLogIdKey() ContextWithValueKey {
	return ContextWithValueKey("LogId")
}

// BrowserFingerprint: a struct to keep fingerprint for client browsers that can be added and passed through context
type BrowserFingerprint struct {
	ClientIP  string
	UserAgent string
}

type LogId struct {
	LogId string
}

// HashPassword : Hash password using salt
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		log.Fatal("unable to hash password", err)
	}

	return string(hash)
}

// GetPasetoKeys : Generate and get paseto prvate key and public keys from secrets respectively
func GetPasetoKeys() (ed25519.PrivateKey, ed25519.PublicKey) {
	b, err := hex.DecodeString(c.PASETO_PRIVATE_KEY_SECRET)
	privateKey := ed25519.PrivateKey(b)

	if err != nil {
		log.Fatalln(err)
	}

	b, err = hex.DecodeString(c.PASETO_PUBLIC_KEY_SECRET)
	publicKey := ed25519.PublicKey(b)

	if err != nil {
		log.Fatalln(err)
	}

	return privateKey, publicKey
}

// TruncateTime : Truncates time with precision for db storage
func TruncateTime(t time.Time) time.Time {
	return time.Unix(0, t.UnixNano()/1e6*1e6)
}

// NewPaseto : Creates a new instance of Paseto
func NewPaseto() *Paseto {
	return &Paseto{
		v2: paseto.NewV2(),
	}
}

// Sign : Signs a payload using a private key
func (p *Paseto) Sign(payload interface{}) (string, error) {
	priv, _ := GetPasetoKeys()

	token, err := p.v2.Sign(priv, &payload, nil)

	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyRefreshToken : Verifies a paseto token and returns the decoded payload
func (p *Paseto) VerifyRefreshToken(token string) (*vo.RefreshTokenPayload, error) {
	payload := &vo.RefreshTokenPayload{}

	_, pub := GetPasetoKeys()

	err := p.v2.Verify(token, pub, payload, nil)

	if err != nil {
		return nil, err
	}

	return payload, nil
}

// VerifyAccessToken : Verifies a paseto token and returns the decoded payload
func (p *Paseto) VerifyAccessToken(token string) (*vo.AccessTokenPayload, error) {
	payload := &vo.AccessTokenPayload{}

	_, pub := GetPasetoKeys()

	err := p.v2.Verify(token, pub, payload, nil)

	if err != nil {
		return nil, err
	}

	return payload, nil
}

// GenerateUserId : Generates a new custom user id
func GenerateUserId() string {
	return fmt.Sprintf("user_%s", generateUUID())
}

func GenerateOrderId() string {
	return fmt.Sprintf("order_%s", generateUUID())
}

// GenerateUserSessionId : Generates a new custom session id
func GenerateUserSessionId() string {
	return fmt.Sprintf("sess_%s", generateUUID())
}

// GenerateOutboxId : Generates a new custom outbox id
func GenerateOutboxId() string {
	return fmt.Sprintf("outb_%s", generateUUID())
}

// GenerateEventId : Generate a new custom event id for identifying events
func GenerateEventId() string {
	return fmt.Sprintf("event_%s", generateUUID())
}

// generateUUID : Generates a new uuid
func generateUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// SetBrowserFingerprintInContext : Sets the browser's ip and agent to the current context
func SetBrowserFingerprintInContext(ctx context.Context, c *fiber.Ctx) context.Context {
	bf := BrowserFingerprint{
		ClientIP:  c.Context().RemoteIP().String(),
		UserAgent: string(c.Context().UserAgent()),
	}

	key := GetBrowserFingerPrintKey()

	return context.WithValue(ctx, key, bf)
}

// GetBrowserFingerPrintFromContext : Gets the browser's ip and agent from the current context
func GetBrowserFingerPrintFromContext(ctx context.Context) BrowserFingerprint {
	key := GetBrowserFingerPrintKey()
	return ctx.Value(key).(BrowserFingerprint)
}

func generateLogId() string {
	return fmt.Sprintf("log_%s", generateUUID())
}

// SetLogIdInContext : Sets the log id to the current context
func SetLogIdInContext(ctx context.Context) (string, context.Context) {
	logId := generateLogId()

	key := GetLogIdKey()

	return logId, context.WithValue(ctx, key, logId)
}

func GetLogIdFromContext(ctx context.Context) string {
	key := GetLogIdKey()
	return ctx.Value(key).(string)
}
