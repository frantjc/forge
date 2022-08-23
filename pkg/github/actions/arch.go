package actions

type Arch string

const (
	ArchX86   Arch = "X86"
	ArchX64   Arch = "X64"
	ArchARM   Arch = "ARM"
	ArchARM64 Arch = "ARM64"
)

func (a Arch) String() string {
	return string(a)
}
