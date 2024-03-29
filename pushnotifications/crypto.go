package pushnotifications

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
)

func (c *Config) GenerateNewKeys() error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	shared := make([]byte, 16)
	_, err = rand.Read(shared)
	if err != nil {
		return err
	}

	c.PrivateKey = *priv
	c.SharedSecret = shared
	return nil
}

func (c *Config) ImportPrivateKey(privateKeyStr string) error {
	k, err := base64.RawURLEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return err
	}

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()

	d := new(big.Int)
	priv.D = d.SetBytes(k)

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(k)
	c.PrivateKey = *priv
	return nil
}

func (c *Config) ImportServerKey(serverKeyStr string) error {
	k, err := base64.RawURLEncoding.DecodeString(serverKeyStr)
	if err != nil {
		return err
	}

	pub := new(ecdsa.PublicKey)
	pub.Curve = elliptic.P256()

	pub.X, pub.Y = elliptic.Unmarshal(pub.Curve, k)
	if pub.X == nil {
		return fmt.Errorf("can not unmarshal server key")
	}

	c.ServerKey = *pub
	return nil
}

func (c *Config) ImportSharedSecret(sharedSecretStr string) error {
	s, err := base64.RawURLEncoding.DecodeString(sharedSecretStr)
	if err != nil {
		return err
	}

	c.SharedSecret = s
	return nil
}

func (c *Config) ExportPrivateKey() string {
	return base64.RawURLEncoding.EncodeToString(c.PrivateKey.D.Bytes())
}

func (c *Config) ExportServerKey() string {
	return base64.RawURLEncoding.EncodeToString(elliptic.Marshal(c.ServerKey.Curve, c.ServerKey.X, c.ServerKey.Y))
}

func (c *Config) ExportSharedSecret() string {
	return base64.RawURLEncoding.EncodeToString(c.SharedSecret)
}

func (m *PushClient) ExportPublicKey() string {
	publicKeyBytes := elliptic.Marshal(m.PrivateKey.PublicKey.Curve, m.PrivateKey.PublicKey.X, m.PrivateKey.PublicKey.Y)
	return base64.RawURLEncoding.EncodeToString(publicKeyBytes)
}
func (m *PushClient) GetPublicKey() (*ecdsa.PublicKey, error) {
	publicKeyBytes := elliptic.Marshal(m.PrivateKey.PublicKey.Curve, m.PrivateKey.PublicKey.X, m.PrivateKey.PublicKey.Y)

	pub := new(ecdsa.PublicKey)
	pub.Curve = elliptic.P256()

	pub.X, pub.Y = elliptic.Unmarshal(pub.Curve, publicKeyBytes)
	if pub.X == nil {
		return nil, fmt.Errorf("can not unmarshal server key")
	}

	return pub, nil
}

func (m *PushClient) deriveSecret(dh []byte) (key []byte, context []byte, err error) {
	x, y := elliptic.Unmarshal(m.PrivateKey.Curve, dh)
	if x == nil {
		return nil, nil, fmt.Errorf("can not unmarshal dh")
	}

	x, _ = m.PrivateKey.Curve.ScalarMult(x, y, m.PrivateKey.D.Bytes())
	return x.Bytes(), m.createContext(dh), nil
}

func (m *PushClient) deriveKey(dh, salt []byte) (key, nonce []byte, err error) {
	hash := sha256.New

	secret, context, err := m.deriveSecret(dh)
	if err != nil {
		return nil, nil, err
	}

	hkdfAuth := NewHkdf(hash, secret, m.SharedSecret, m.createInfo("auth", nil))

	newSecret := make([]byte, 32)
	if _, err := io.ReadFull(hkdfAuth, newSecret); err != nil {
		return nil, nil, err
	}

	hkdfKey := NewHkdf(hash, newSecret, salt, m.createInfo("aesgcm", context))  // Length: 16
	hkdfNonce := NewHkdf(hash, newSecret, salt, m.createInfo("nonce", context)) // Length: 12

	key = make([]byte, 16)
	nonce = make([]byte, 12)

	if _, err := io.ReadFull(hkdfKey, key); err != nil {
		return nil, nil, err
	}

	if _, err := io.ReadFull(hkdfNonce, nonce); err != nil {
		return nil, nil, err
	}

	return key, nonce, nil
}

func (m *PushClient) createInfo(eType string, context []byte) []byte {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("Content-Encoding: %s\x00", eType))
	buf.Write(context)
	return buf.Bytes()
}

func (m *PushClient) createContext(dh []byte) (context []byte) {
	publicKeyBytes := elliptic.Marshal(m.PrivateKey.Curve, m.PrivateKey.X, m.PrivateKey.Y)

	buf := new(bytes.Buffer)
	buf.WriteString("P-256\x00")

	length := make([]byte, 2)

	binary.BigEndian.PutUint16(length, uint16(len(publicKeyBytes)))
	buf.Write(length)
	buf.Write(publicKeyBytes)

	binary.BigEndian.PutUint16(length, uint16(len(dh)))
	buf.Write(length)
	buf.Write(dh)

	return buf.Bytes()
}

func (m *PushClient) Decrypt(dh, salt, data []byte) (payload []byte, err error) {
	key, nonce, err := m.deriveKey(dh, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	payload, err = aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	return payload[2:], nil
}
