package IP

import "time"

//change detial to details
type IPpost struct {
	IPaddress string    `json:"ip"`
	Detail    IPdetails `json:"detail"`
}

type IPdetails struct {
	MACaddress string    `json:"MACaddress"`
	LeaseTime  time.Time `json:"leaseTime"`
	Available  bool      `json:"available"`
}
