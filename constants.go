package uzimiddleware

// OID number of the UZI card
type UziOidCA string

const (
	OID_CA_CARE_PROVIDER    UziOidCA = "2.16.528.1.1003.1.3.5.5.2" // Reference page 59
	OID_CA_NAMED_EMPLOYEE   UziOidCA = "2.16.528.1.1003.1.3.5.5.3"
	OID_CA_UNNAMED_EMPLOYEE UziOidCA = "2.16.528.1.1003.1.3.5.5.4"
	OID_CA_SERVER           UziOidCA = "2.16.528.1.1003.1.3.5.5.5"
)

// Type of the UZI card
type UziType string

const (
	UZI_TYPE_CARE_PROVIDER    UziType = "Z" // Reference page 60
	UZI_TYPE_NAMED_EMPLOYEE   UziType = "N"
	UZI_TYPE_UNNAMED_EMPLOYEE UziType = "M"
	UZI_TYPE_SERVER           UziType = "S"
)

// Role of the UZI card holder
type UziRole string

const (
	UZI_ROLE_PHARMACIST                UziRole = "17." // Reference page 89
	UZI_ROLE_DOCTOR                    UziRole = "01."
	UZI_ROLE_PHYSIOTHERAPIST           UziRole = "04."
	UZI_ROLE_HEALTHCARE_PSYCHOLOGIST   UziRole = "25."
	UZI_ROLE_PSYCHOTHERAPIST           UziRole = "16."
	UZI_ROLE_DENTIST                   UziRole = "02."
	UZI_ROLE_MIDWIFE                   UziRole = "02."
	UZI_ROLE_NURSE                     UziRole = "30."
	UZI_ROLE_PHYS_ASSISTANT            UziRole = "81."
	UZI_ROLE_ORTHOPEDAGOGUE_GENERALIST UziRole = "31."
	UZI_ROLE_CLINICAL_TECHNOLOGIST     UziRole = "82."
)
