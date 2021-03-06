package jwe

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/lestrrat/go-jwx/buffer"
	"github.com/lestrrat/go-jwx/internal/debug"
	"github.com/lestrrat/go-jwx/internal/emap"
	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwk"
)

func NewRecipient() *Recipient {
	return &Recipient{
		Header: NewHeader(),
	}
}

func NewHeader() *Header {
	return &Header{
		EssentialHeader: &EssentialHeader{},
		PrivateParams:   map[string]interface{}{},
	}
}

func NewEncodedHeader() *EncodedHeader {
	return &EncodedHeader{
		Header: NewHeader(),
	}
}

func (h *Header) Get(key string) (interface{}, error) {
	switch key {
	case "alg":
		return h.Algorithm, nil
	case "apu":
		return h.AgreementPartyUInfo, nil
	case "apv":
		return h.AgreementPartyVInfo, nil
	case "enc":
		return h.ContentEncryption, nil
	case "epk":
		return h.EphemeralPublicKey, nil
	case "cty":
		return h.ContentType, nil
	case "kid":
		return h.KeyID, nil
	case "typ":
		return h.Type, nil
	case "x5t":
		return h.X509CertThumbprint, nil
	case "x5t#256":
		return h.X509CertThumbprintS256, nil
	case "x5c":
		return h.X509CertChain, nil
	case "crit":
		return h.Critical, nil
	case "jku":
		return h.JwkSetURL, nil
	case "x5u":
		return h.X509Url, nil
	default:
		v, ok := h.PrivateParams[key]
		if !ok {
			return nil, errors.New("invalid header name")
		}
		return v, nil
	}
}

func (h *Header) Set(key string, value interface{}) error {
	switch key {
	case "alg":
		var v jwa.KeyEncryptionAlgorithm
		s, ok := value.(string)
		if ok {
			v = jwa.KeyEncryptionAlgorithm(s)
		} else {
			v, ok = value.(jwa.KeyEncryptionAlgorithm)
			if !ok {
				return ErrInvalidHeaderValue
			}
		}
		h.Algorithm = v
	case "apu":
		var v buffer.Buffer
		switch value.(type) {
		case buffer.Buffer:
			v = value.(buffer.Buffer)
		case []byte:
			v = buffer.Buffer(value.([]byte))
		case string:
			v = buffer.Buffer(value.(string))
		default:
			return ErrInvalidHeaderValue
		}
		h.AgreementPartyUInfo = v
	case "apv":
		var v buffer.Buffer
		switch value.(type) {
		case buffer.Buffer:
			v = value.(buffer.Buffer)
		case []byte:
			v = buffer.Buffer(value.([]byte))
		case string:
			v = buffer.Buffer(value.(string))
		default:
			return ErrInvalidHeaderValue
		}
		h.AgreementPartyVInfo = v
	case "enc":
		var v jwa.ContentEncryptionAlgorithm
		s, ok := value.(string)
		if ok {
			v = jwa.ContentEncryptionAlgorithm(s)
		} else {
			v, ok = value.(jwa.ContentEncryptionAlgorithm)
			if !ok {
				return ErrInvalidHeaderValue
			}
		}
		h.ContentEncryption = v
	case "cty":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.ContentType = v
	case "epk":
		v, ok := value.(*jwk.EcdsaPublicKey)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.EphemeralPublicKey = v
	case "kid":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.KeyID = v
	case "typ":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.Type = v
	case "x5t":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.X509CertThumbprint = v
	case "x5t#256":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.X509CertThumbprintS256 = v
	case "x5c":
		v, ok := value.([]string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.X509CertChain = v
	case "crit":
		v, ok := value.([]string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.Critical = v
	case "jku":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		u, err := url.Parse(v)
		if err != nil {
			return ErrInvalidHeaderValue
		}
		h.JwkSetURL = u
	case "x5u":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		u, err := url.Parse(v)
		if err != nil {
			return ErrInvalidHeaderValue
		}
		h.X509Url = u
	default:
		h.PrivateParams[key] = value
	}
	return nil
}

func (h1 *Header) Merge(h2 *Header) (*Header, error) {
	if h2 == nil {
		return nil, errors.New("merge target is nil")
	}

	h3 := NewHeader()
	if err := h3.Copy(h1); err != nil {
		return nil, err
	}

	h3.EssentialHeader.Merge(h2.EssentialHeader)

	for k, v := range h2.PrivateParams {
		h3.PrivateParams[k] = v
	}

	return h3, nil
}

func (h1 *EssentialHeader) Merge(h2 *EssentialHeader) {
	if h2.AgreementPartyUInfo.Len() != 0 {
		h1.AgreementPartyUInfo = h2.AgreementPartyUInfo
	}

	if h2.AgreementPartyVInfo.Len() != 0 {
		h1.AgreementPartyVInfo = h2.AgreementPartyVInfo
	}

	if h2.Algorithm != "" {
		h1.Algorithm = h2.Algorithm
	}

	if h2.ContentEncryption != "" {
		h1.ContentEncryption = h2.ContentEncryption
	}

	if h2.ContentType != "" {
		h1.ContentType = h2.ContentType
	}

	if h2.Compression != "" {
		h1.Compression = h2.Compression
	}

	if h2.Critical != nil {
		h1.Critical = h2.Critical
	}

	if h2.EphemeralPublicKey != nil {
		h1.EphemeralPublicKey = h2.EphemeralPublicKey
	}

	if h2.Jwk != nil {
		h1.Jwk = h2.Jwk
	}

	if h2.JwkSetURL != nil {
		h1.JwkSetURL = h2.JwkSetURL
	}

	if h2.KeyID != "" {
		h1.KeyID = h2.KeyID
	}

	if h2.Type != "" {
		h1.Type = h2.Type
	}

	if h2.X509Url != nil {
		h1.X509Url = h2.X509Url
	}

	if h2.X509CertChain != nil {
		h1.X509CertChain = h2.X509CertChain
	}

	if h2.X509CertThumbprint != "" {
		h1.X509CertThumbprint = h2.X509CertThumbprint
	}

	if h2.X509CertThumbprintS256 != "" {
		h1.X509CertThumbprintS256 = h2.X509CertThumbprintS256
	}
}

func (h1 *Header) Copy(h2 *Header) error {
	if h1 == nil {
		return errors.New("copy destination is nil")
	}
	if h2 == nil {
		return errors.New("copy target is nil")
	}

	h1.EssentialHeader.Copy(h2.EssentialHeader)

	for k, v := range h2.PrivateParams {
		h1.PrivateParams[k] = v
	}

	return nil
}

func (h1 *EssentialHeader) Copy(h2 *EssentialHeader) {
	h1.AgreementPartyUInfo = h2.AgreementPartyUInfo
	h1.AgreementPartyVInfo = h2.AgreementPartyVInfo
	h1.Algorithm = h2.Algorithm
	h1.ContentEncryption = h2.ContentEncryption
	h1.ContentType = h2.ContentType
	h1.Compression = h2.Compression
	h1.Critical = h2.Critical
	h1.EphemeralPublicKey = h2.EphemeralPublicKey
	h1.Jwk = h2.Jwk
	h1.JwkSetURL = h2.JwkSetURL
	h1.KeyID = h2.KeyID
	h1.Type = h2.Type
	h1.X509Url = h2.X509Url
	h1.X509CertChain = h2.X509CertChain
	h1.X509CertThumbprint = h2.X509CertThumbprint
	h1.X509CertThumbprintS256 = h2.X509CertThumbprintS256
}

func (h Header) MarshalJSON() ([]byte, error) {
	return emap.MergeMarshal(h.EssentialHeader, h.PrivateParams)
}

func (h *Header) UnmarshalJSON(data []byte) error {
	if h.EssentialHeader == nil {
		h.EssentialHeader = &EssentialHeader{}
	}
	if h.PrivateParams == nil {
		h.PrivateParams = map[string]interface{}{}
	}

	if err := json.Unmarshal(data, h.EssentialHeader); err != nil {
		return err
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for _, n := range []string{"alg", "apu", "apv", "enc", "cty", "zip", "crit", "epk", "jwk", "jku", "kid", "typ", "x5u", "x5c", "x5t", "x5t#S256"} {
		delete(m, n)
	}

	for name, value := range m {
		if err := h.Set(name, value); err != nil {
			return err
		}
	}
	return nil
}

func (e EncodedHeader) Base64Encode() ([]byte, error) {
	buf, err := json.Marshal(e.Header)
	if err != nil {
		return nil, err
	}

	buf, err = buffer.Buffer(buf).Base64Encode()
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (e EncodedHeader) MarshalJSON() ([]byte, error) {
	buf, err := e.Base64Encode()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(buf))
}

func (e *EncodedHeader) UnmarshalJSON(buf []byte) error {
	b := buffer.Buffer{}
	// base646 json string -> json object representation of header
	if err := json.Unmarshal(buf, &b); err != nil {
		return err
	}

	if err := json.Unmarshal(b.Bytes(), &e.Header); err != nil {
		return err
	}

	return nil
}

func NewMessage() *Message {
	return &Message{
		ProtectedHeader:   NewEncodedHeader(),
		UnprotectedHeader: NewHeader(),
	}
}

func (m *Message) Decrypt(alg jwa.KeyEncryptionAlgorithm, key interface{}) ([]byte, error) {
	var err error

	if len(m.Recipients) == 0 {
		return nil, errors.New("no recipients, can not proceed with decrypt")
	}

	enc := m.ProtectedHeader.ContentEncryption

	h := NewHeader()
	if err := h.Copy(m.ProtectedHeader.Header); err != nil {
		return nil, err
	}
	h, err = h.Merge(m.UnprotectedHeader)
	if err != nil {
		debug.Printf("failed to merge unprotected header")
		return nil, err
	}

	aad, err := m.AuthenticatedData.Base64Encode()
	if err != nil {
		return nil, err
	}
	ciphertext := m.CipherText.Bytes()
	iv := m.InitializationVector.Bytes()
	tag := m.Tag.Bytes()

	cipher, err := BuildContentCipher(enc)
	if err != nil {
		return nil, fmt.Errorf("unsupported content cipher algorithm '%s'", enc)
	}
	keysize := cipher.KeySize()

	var plaintext []byte
	for _, recipient := range m.Recipients {
		debug.Printf("Attempting to check if we can decode for recipient (alg = %s)", recipient.Header.Algorithm)
		if recipient.Header.Algorithm != alg {
			continue
		}

		h2 := NewHeader()
		if err := h2.Copy(h); err != nil {
			debug.Printf("failed to copy header: %s", err)
			continue
		}

		h2, err := h2.Merge(recipient.Header)
		if err != nil {
			debug.Printf("Failed to merge! %s", err)
			continue
		}

		k, err := BuildKeyDecrypter(h2.Algorithm, h2, key, keysize)
		if err != nil {
			debug.Printf("failed to create key decrypter: %s", err)
			continue
		}

		cek, err := k.KeyDecrypt(recipient.EncryptedKey.Bytes())
		if err != nil {
			debug.Printf("failed to decrypt key: %s", err)
			return nil, errors.New("failed to decrypt key")
			continue
		}

		plaintext, err = cipher.decrypt(cek, iv, ciphertext, tag, aad)
		if err == nil {
			break
		}
		debug.Printf("DecryptMessage: failed to decrypt using %s: %s", h2.Algorithm, err)
		// Keep looping because there might be another key with the same algo
	}

	if plaintext == nil {
		return nil, errors.New("failed to find matching recipient to decrypt key")
	}

	if h.Compression == jwa.Deflate {
		output := bytes.Buffer{}
		w, _ := flate.NewWriter(&output, 1)
		in := plaintext
		for len(in) > 0 {
			n, err := w.Write(in)
			if err != nil {
				return nil, err
			}
			in = in[n:]
		}
		if err := w.Close(); err != nil {
			return nil, err
		}
		plaintext = output.Bytes()
	}

	return plaintext, nil
}
