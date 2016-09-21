package mongodb

type Config struct {
	Init     bool
	Hostname string `valid:"host"`
	Port     string `valid:"port"`
	Username string
	Password string
	Database string
	// Direct if enabled will disables the automatic replica set server discovery logic, and
	// forces the use of servers provided only (even if secondaries).
	// Note that to talk to a secondary the consistency requirements
	// must be relaxed to Monotonic or Eventual via SetMode.
	Direct bool
}
