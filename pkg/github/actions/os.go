package actions

type OS string

const (
	OSLinux   OS = "Linux"
	OSWindows OS = "Windows"
	OSDarwin  OS = "macOS"
)

func (o OS) String() string {
	return string(o)
}
