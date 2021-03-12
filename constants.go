package uzimiddleware

// UziOidCA is the OID number of the UZI card
type UziOidCA string

// Predefined OIC CAs
const (
	OidCaCareProvider    UziOidCA = "2.16.528.1.1003.1.3.5.5.2" // Reference page 59
	OidCaNamedEmployee   UziOidCA = "2.16.528.1.1003.1.3.5.5.3"
	OidCaUnnamedEmployee UziOidCA = "2.16.528.1.1003.1.3.5.5.4"
	OidCaServer          UziOidCA = "2.16.528.1.1003.1.3.5.5.5"
)

// UziType is the type of the UZI card
type UziType string

// Predefined UZI types
const (
	UziTypeCareProvider    UziType = "Z" // Reference page 60
	UziTypeNamedEmployee   UziType = "N"
	UziTypeUnnamedEmployee UziType = "M"
	UziTypeServer          UziType = "S"
)

// UziRole is the role of the UZI card holder
type UziRole string

// Predefined UZI roles
const (
	UziRolePharmacist               UziRole = "17." // Reference page 89
	UziRoleDoctor                   UziRole = "01."
	UziRolePhysiotherapist          UziRole = "04."
	UziRoleHealthcarePsychologist   UziRole = "25."
	UziRolePsychotherapist          UziRole = "16."
	UziRoleDentist                  UziRole = "02."
	UziRoleMidwife                  UziRole = "03."
	UziRoleNurse                    UziRole = "30."
	UziRolePhysAssistant            UziRole = "81."
	UziRoleOrthopedagogueGeneralist UziRole = "31."
	UziRoleClinicalTechnologist     UziRole = "82."
)
