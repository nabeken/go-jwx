package jws

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/lestrrat/go-jwx/buffer"
	"github.com/lestrrat/go-jwx/internal/emap"
	"github.com/lestrrat/go-jwx/jwa"
)

func NewHeader() *Header {
	return &Header{
		EssentialHeader: &EssentialHeader{},
		PrivateParams:   map[string]interface{}{},
	}
}

func (h *Header) Set(key string, value interface{}) error {
	switch key {
	case "alg":
		var v jwa.SignatureAlgorithm
		s, ok := value.(string)
		if ok {
			v = jwa.SignatureAlgorithm(s)
		} else {
			v, ok = value.(jwa.SignatureAlgorithm)
			if !ok {
				return ErrInvalidHeaderValue
			}
		}
		h.Algorithm = v
	case "cty":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.ContentType = v
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
	if h2.Algorithm != "" {
		h1.Algorithm = h2.Algorithm
	}

	if h2.ContentType != "" {
		h1.ContentType = h2.ContentType
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
  h1.Algorithm = h2.Algorithm
  h1.ContentType = h2.ContentType
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
	return emap.MergeUnmarshal(data, h.EssentialHeader, &h.PrivateParams)
}

func (h *EssentialHeader) Construct(m map[string]interface{}) error {
	r := emap.Hmap(m)
	if alg, err := r.GetString("alg"); err == nil {
		h.Algorithm = jwa.SignatureAlgorithm(alg)
	}
	if h.Algorithm == "" {
		h.Algorithm = jwa.NoSignature
	}
	h.ContentType, _ = r.GetString("cty")
	h.KeyID, _ = r.GetString("kid")
	h.Type, _ = r.GetString("typ")
	h.X509CertThumbprint, _ = r.GetString("x5t")
	h.X509CertThumbprintS256, _ = r.GetString("x5t#256")
	if v, err := r.GetStringSlice("crit"); err != nil {
		h.Critical = v
	}
	if v, err := r.GetStringSlice("x5c"); err != nil {
		h.X509CertChain = v
	}
	if v, err := r.GetString("jku"); err == nil {
		u, err := url.Parse(v)
		if err == nil {
			h.JwkSetURL = u
		}
	}

	if v, err := r.GetString("x5u"); err == nil {
		u, err := url.Parse(v)
		if err == nil {
			h.X509Url = u
		}
	}

	return nil
}

func (h Header) Base64Encode() ([]byte, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}

	return buffer.Buffer(b).Base64Encode()
}

func (e EncodedHeader) MarshalJSON() ([]byte, error) {
	buf, err := json.Marshal(e.Header)
	if err != nil {
		return nil, err
	}

	buf, err = buffer.Buffer(buf).Base64Encode()
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

	e.Source = b

	return nil
}

func NewSignature() *Signature {
	h1 := NewHeader()
	h2 := NewHeader()
	return &Signature{
		PublicHeader:    h1,
		ProtectedHeader: &EncodedHeader{Header: h2},
	}
}

func (s Signature) MergedHeaders() MergedHeader {
	return MergedHeader{
		ProtectedHeader: s.ProtectedHeader,
		PublicHeader:    s.PublicHeader,
	}
}

func (h MergedHeader) KeyID() string {
	if hp := h.ProtectedHeader; hp != nil {
		if hp.KeyID != "" {
			return hp.KeyID
		}
	}

	if hp := h.PublicHeader; hp != nil {
		if hp.KeyID != "" {
			return hp.KeyID
		}
	}

	return ""
}

func (h MergedHeader) Algorithm() jwa.SignatureAlgorithm {
	if hp := h.ProtectedHeader; hp != nil {
		return hp.Algorithm
	}
	return jwa.NoSignature
}

// LookupSignature looks up a particular signature entry using
// the `kid` value
func (m Message) LookupSignature(kid string) []Signature {
	sigs := []Signature{}
	for _, sig := range m.Signatures {
		if sig.MergedHeaders().KeyID() != kid {
			continue
		}

		sigs = append(sigs, sig)
	}
	return sigs
}
