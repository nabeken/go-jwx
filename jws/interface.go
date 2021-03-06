package jws

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"hash"
	"net/url"

	"github.com/lestrrat/go-jwx/buffer"
	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwk"
)

var (
	ErrInvalidCompactPartsCount  = errors.New("compact JWS format must have three parts")
	ErrInvalidHeaderValue        = errors.New("invalid value for header key")
	ErrInvalidEcdsaSignatureSize = errors.New("invalid signature size of ecdsa algorithm")
	ErrInvalidSignature          = errors.New("invalid signature")
	ErrMissingPrivateKey         = errors.New("missing private key")
	ErrMissingPublicKey          = errors.New("missing public key")
	ErrUnsupportedAlgorithm      = errors.New("unspported algorithm")
)

type EssentialHeader struct {
	Algorithm              jwa.SignatureAlgorithm `json:"alg,omitempty"`
	ContentType            string                 `json:"cty,omitempty"`
	Critical               []string               `json:"crit,omitempty"`
	Jwk                    jwk.Key                `json:"jwk,omitempty"` // public key
	JwkSetURL              *url.URL               `json:"jku,omitempty"`
	KeyID                  string                 `json:"kid,omitempty"`
	Type                   string                 `json:"typ,omitempty"` // e.g. "JWT"
	X509Url                *url.URL               `json:"x5u,omitempty"`
	X509CertChain          []string               `json:"x5c,omitempty"`
	X509CertThumbprint     string                 `json:"x5t,omitempty"`
	X509CertThumbprintS256 string                 `json:"x5t#S256,omitempty"`
}

// Header represents a jws header.
type Header struct {
	*EssentialHeader `json:"-"`
	PrivateParams    map[string]interface{} `json:"-"`
}

// EncodedHeader represents a header value that is base64 encoded
// in JSON format
type EncodedHeader struct {
	*Header
	// This is a special field. It's ONLY set when parsed from a serialized form.
	// It's used for verification purposes, because header representations (such as
	// JSON key order) may differ from what the source encoded with and what the
	// go json package uses
	//
	// If this field is populated (Len() > 0), it will be used for signature
	// calculation.
	// If you change the header values, make sure to clear this field, too
	Source buffer.Buffer `json:"-"`
}

// PayloadSigner generates signature for the given payload
type PayloadSigner interface {
	PayloadSign([]byte) ([]byte, error)
	PublicHeaders() *Header
	ProtectedHeaders() *Header
	SetPublicHeaders(*Header)
	SetProtectedHeaders(*Header)
	SignatureAlgorithm() jwa.SignatureAlgorithm
}

// Verifier is used to verify the signature against the payload
type Verifier interface {
	Verify(*Message) error
}

type RsaSign struct {
	Public     *Header
	Protected  *Header
	PrivateKey *rsa.PrivateKey
}

type EcdsaSign struct {
	Public     *Header
	Protected  *Header
	PrivateKey *ecdsa.PrivateKey
}

type MergedHeader struct {
	ProtectedHeader *EncodedHeader
	PublicHeader    *Header
}

type Signature struct {
	PublicHeader    *Header        `json:"header"`              // Raw JWS Unprotected Heders
	ProtectedHeader *EncodedHeader `json:"protected,omitempty"` // Base64 encoded JWS Protected Headers
	Signature       buffer.Buffer  `json:"signature"`           // Base64 encoded signature
}

// Message represents a full JWS encoded message. Flattened serialization
// is not supported as a struct, but rather it's represented as a
// Message struct with only one `signature` element
type Message struct {
	Payload    buffer.Buffer `json:"payload"`
	Signatures []Signature   `json:"signatures"`
}

type MultiSign struct {
	Signers []PayloadSigner
}

type HmacSign struct {
	Public     *Header
	Protected  *Header
	Key       []byte
	hash      func() hash.Hash
}

type Serializer interface {
	Serialize(*Message) ([]byte, error)
}

type CompactSerialize struct{}
type JSONSerialize struct {
	Pretty bool
}

type RsaVerify struct {
	alg    jwa.SignatureAlgorithm
	hash   crypto.Hash
	pubkey *rsa.PublicKey
}

type EcdsaVerify struct {
	alg    jwa.SignatureAlgorithm
	hash   crypto.Hash
	pubkey *ecdsa.PublicKey
}
