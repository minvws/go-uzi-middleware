package uzimiddleware

import (
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"strings"
)

var (
	errInvalidSubjectAltName = errors.New("uzi: invalid SAN specified")
	errIncorrectOID          = errors.New("uzi: CA OID not UZI register Care Provider or named employee")
	errIncorrectUziVersion   = errors.New("uzi: Version not 1")
	errCardTypeNotAllowed    = errors.New("uzi: card type not allowed")
	errCardRoleNotAllowed    = errors.New("uzi: role not allowed")
)

// UziUser is the populated structure that is returned by context to the request. It contains all information about the found UZI card
type UziUser struct {
	AgbCode          string   // AGB code
	CardType         UziType  // Type of the card
	GivenName        string   // Given name of the UZI card holder
	OidCA            UziOidCA // OID CA of the card
	Role             UziRole  // Role of the card holder
	SubscriberNumber string   // Subscriber number of the UZI card
	SurName          string   // Surname of the UZI card holder
	UziNumber        string   // Number of the UZI card
	UziVersion       string   // Version of the UZI card
	UPN              string   // Microsoft UPN number
}

type value struct {
	S string `asn1:"optional,omitempty"`
}

type otherName struct {
	OID   asn1.ObjectIdentifier
	Value value `asn1:"tag:0"`
}

var (
	idAtSurName        = asn1.ObjectIdentifier{2, 5, 4, 4}
	idAtGivenName      = asn1.ObjectIdentifier{2, 5, 4, 42}
	idCeSubjectAltName = asn1.ObjectIdentifier{2, 5, 29, 17}

	idIA5string = asn1.ObjectIdentifier{2, 5, 5, 5}
	idUPN       = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 20, 2, 3}
)

// NewUziUserFromCert populates a UZI user structure based on the data found in the certificate, or returns an error
func NewUziUserFromCert(cert *x509.Certificate) (*UziUser, error) {
	u := &UziUser{}

	// Extract surname and given name from subject
	for _, attr := range cert.Subject.Names {
		if attr.Type.Equal(idAtSurName) {
			u.SurName = attr.Value.(string)
		}
		if attr.Type.Equal(idAtGivenName) {
			u.GivenName = attr.Value.(string)
		}
	}

	// extract data from SubjectAltName extension
	for i := range cert.Extensions {
		ext := cert.Extensions[i]
		if !ext.Id.Equal(idCeSubjectAltName) {
			continue
		}

		otherNames, err := parseOtherNames(ext.Value)
		if err != nil {
			return nil, err
		}

		for _, on := range otherNames {
			// UPN
			if on.OID.Equal(idUPN) {
				u.UPN = on.Value.S
			}

			// Regular string
			if on.OID.Equal(idIA5string) {
				parts := strings.Split(on.Value.S, "-")
				if len(parts) >= 6 {
					u.OidCA = UziOidCA(parts[0])
					u.UziVersion = parts[1]
					u.UziNumber = parts[2]
					u.CardType = UziType(parts[3])
					u.SubscriberNumber = parts[4]
					u.Role = UziRole(parts[5])
					u.AgbCode = parts[6]
				}
			}
		}

		// Check if data is filled in. If so, we are done
		if u.SubscriberNumber != "" && u.UziNumber != "" {
			return u, nil
		}
	}

	return nil, errInvalidSubjectAltName
}

func parseOtherNames(bytes []byte) ([]otherName, error) {
	var seq asn1.RawValue
	_, err := asn1.Unmarshal(bytes, &seq)
	if err != nil {
		return nil, errInvalidSubjectAltName
	}

	if !seq.IsCompound || seq.Tag != asn1.TagSequence || seq.Class != asn1.ClassUniversal {
		return nil, errInvalidSubjectAltName
	}

	var otherNames []otherName

	rest := seq.Bytes
	for len(rest) > 0 {
		var err error
		otherNames, rest, err = parseOtherName(rest, otherNames)
		if err != nil {
			return nil, errInvalidSubjectAltName
		}
	}

	return otherNames, nil
}

func parseOtherName(bytes []byte, othernames []otherName) ([]otherName, []byte, error) {
	var on otherName

	bytes = append([]byte{}, bytes...)
	bytes[0] = asn1.TagSequence | 0x20
	rest, err := asn1.Unmarshal(bytes, &on)
	if err != nil {
		return othernames, nil, err
	}

	return append(othernames, on), rest, nil
}

// Validate will validate an UZI user against the given arguments. Will return nil when validated correctly, or error
func (u *UziUser) Validate(strictCA bool, allowedTypes []UziType, allowedRoles []UziRole) error {
	if strictCA && u.OidCA != OidCaCareProvider && u.OidCA != OidCaNamedEmployee {
		return errIncorrectOID
	}

	if u.UziVersion != "1" {
		return errIncorrectUziVersion
	}

	if !containsType(allowedTypes, u.CardType) {
		return errCardTypeNotAllowed
	}

	if !containsRole(allowedRoles, u.Role) {
		return errCardRoleNotAllowed
	}

	return nil
}

func containsType(s []UziType, e UziType) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsRole(s []UziRole, e UziRole) bool {
	for _, a := range s {
		if a[:3] == e[:3] {
			return true
		}
	}
	return false
}
