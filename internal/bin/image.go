package bin

const (
	ShimPath = "/" + ShimName
)

var (
	ShimEntrypoint      = []string{ShimPath}
	ShimSleepEntrypoint = append(ShimEntrypoint, "-s")
)
