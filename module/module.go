package module

type Module struct {
	Inspeed     uint64
	Outspeed    uint64
	Connection  int32
	Memoryusage uint64
}
type Flag struct {
	Addr string `short:"a" long:"addr" description:"Web server address REQUIRED"  `
	NIC  int    `short:"n" long:"nic" description:"NIC number REQUIRED" `
	List bool   `short:"l" long:"list" description:"List all NICs"`
}
