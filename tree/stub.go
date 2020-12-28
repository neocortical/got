package tree

type stubNode struct {
	name string
	oid  string
	mode string
}

func (sn stubNode) ModeString() string {
	return sn.mode
}

func (sn stubNode) Name() string {
	return sn.name
}

func (sn stubNode) OID() string {
	return sn.oid
}
