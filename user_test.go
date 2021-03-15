package uzimiddleware

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUziUserFromCert(t *testing.T) {
	var (
		u    *UziUser
		err  error
		cert *x509.Certificate
	)

	cert = loadCert(t, "./testdata/mock-001-no-rnd.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-001-no-valid-uzi-data.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-002.crt")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-002-invalid-san.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-003.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-003-invalid-othername.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-003-valid-san.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-004-othername-without-ia5string.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-004-valid-othername.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-005-incorrect-san-data.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-005-valid-othername.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-006-incorrect-san-data.cert")
	u, err = NewUziUserFromCert(cert)
	assert.ErrorIs(t, err, errInvalidSubjectAltName)
	assert.Nil(t, u)

	cert = loadCert(t, "./testdata/mock-011-correct.cert")
	u, err = NewUziUserFromCert(cert)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, "12345678", u.UziNumber)
	assert.Equal(t, UziTypeNamedEmployee, u.CardType)
	assert.Equal(t, "90000111", u.SubscriberNumber)
	assert.Equal(t, "john", u.GivenName)
	assert.Equal(t, "doe-12345678", u.SurName)

	cert = loadCert(t, "./testdata/mock-012-correct-admin.cert")
	u, err = NewUziUserFromCert(cert)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, "11111111", u.UziNumber)
	assert.Equal(t, UziTypeNamedEmployee, u.CardType)
	assert.Equal(t, "90000111", u.SubscriberNumber)
	assert.Equal(t, "john", u.GivenName)
	assert.Equal(t, "doe-11111111", u.SurName)
}

func TestUziUser_Validate(t *testing.T) {
	roles := []UziRole{UziRoleDoctor, UziRoleDentist}
	types := []UziType{UziTypeNamedEmployee, UziTypeCareProvider}

	cert := loadCert(t, "./testdata/mock-012-correct-admin.cert")
	u, _ := NewUziUserFromCert(cert)
	assert.Error(t, u.Validate(true, nil, nil))
	assert.NoError(t, u.Validate(true, types, roles))

	cert = loadCert(t, "./testdata/mock-007-strict-ca-check.cert")
	u, _ = NewUziUserFromCert(cert)
	assert.Error(t, u.Validate(true, types, roles))
	assert.NoError(t, u.Validate(false, types, roles))

	cert = loadCert(t, "./testdata/mock-008-invalid-version.cert")
	u, _ = NewUziUserFromCert(cert)
	err := u.Validate(true, types, roles)
	assert.ErrorIs(t, err, errIncorrectUziVersion)

	cert = loadCert(t, "./testdata/mock-009-invalid-types.cert")
	u, _ = NewUziUserFromCert(cert)
	err = u.Validate(true, types, roles)
	assert.ErrorIs(t, err, errCardTypeNotAllowed)

	cert = loadCert(t, "./testdata/mock-010-invalid-roles.cert")
	u, _ = NewUziUserFromCert(cert)
	err = u.Validate(true, types, roles)
	assert.ErrorIs(t, err, errCardRoleNotAllowed)
}

func loadCert(t *testing.T, certFile string) *x509.Certificate {
	buf, err := ioutil.ReadFile(certFile)
	if err != nil {
		t.Fatal("cannot load ", certFile)
	}

	block, _ := pem.Decode([]byte(buf))
	if block == nil {
		t.Fatal("cannot decode ", certFile)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal("cannot parse ", certFile, ": ", err)
	}

	return cert
}
