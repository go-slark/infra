package uid

import (
	"github.com/bwmarrin/snowflake"
	"github.com/rs/xid"
	"net"
	"os"
)

// snowflake (8bytes)

type Node struct {
	*snowflake.Node
}

// nodeID : 0 --> 1023

func NewNode(nodeID ...int64) (*Node, error) {
	var id int64
	if len(nodeID) == 0 {
		cid, err := centerID()
		if err != nil {
			return nil, err
		}
		wid, err := workID()
		if err != nil {
			return nil, err
		}
		id = cid | wid
	} else {
		id = nodeID[0]
	}
	node, err := snowflake.NewNode(id)
	if err != nil {
		return nil, err
	}
	return &Node{Node: node}, nil
}

func (n *Node) GenerateID() int64 {
	return n.Generate().Int64()
}

// center id 5bits
func centerID() (int64, error) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return 0, err
	}

	var (
		in *net.IPNet
		ok bool
		ip string
	)
	for _, addr := range address {
		in, ok = addr.(*net.IPNet)
		if !ok || in.IP.IsLoopback() {
			continue
		}
		if in.IP.To4() == nil {
			continue
		}
		ip = in.IP.String()
		break
	}
	var sum uint8
	for _, bs := range []byte(ip) {
		sum += bs
	}
	return int64(sum%32) << 5, nil
}

// work id 5bits
func workID() (int64, error) {
	hn, err := os.Hostname()
	if err != nil {
		return 0, nil
	}
	var sum uint8
	for _, bs := range []byte(hn) {
		sum += bs
	}
	return int64(sum % 32), nil
}

// xid (12bytes)

func GenerateID() string {
	return xid.New().String()
}
