package efi

//go:generate go run github.com/hexaflex/stringer -flags -type=Attributes attributes.go

type Attributes uint32

const (
	// Each variable has Attributes that define how the firmware
	// stores and maintains the data value.

	// NonVolatile indicates the firmware environment variable is
	// stored in non-volatile memory (e.g. NVRAM).
	NonVolatile Attributes = 0x0001

	// BootServiceAccess indicates whenever the variable is visible
	// for boot services.
	BootServiceAccess Attributes = 0x0002

	// RuntimeAccess indicates whenever the variable is visible
	// after ExitBootServices() in UEFI is called.
	RuntimeAccess Attributes = 0x0004

	// HardwareErrorRecord marks the variable to be stored in an
	// area for hardware errors.
	HardwareErrorRecord Attributes = 0x0008

	// AuthenticatedWriteAccess is deprecated and should no longer
	// be used.
	AuthenticatedWriteAccess Attributes = 0x0010

	TimeBasedAuthenticatedWriteAccess Attributes = 0x0020

	AppendWrite Attributes = 0x0040

	// EnhancedAuthenticatedAccess attribute indicates that the
	// variable payload begins with an
	// EFI_VARIABLE_AUTHENTICATION_3 structure, and potentially
	// more structures as indicated by fields of this structure.
	EnhancedAuthenticatedAccess Attributes = 0x0080
)
