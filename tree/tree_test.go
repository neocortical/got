package tree

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestDeserializeTree(t *testing.T) {
	data, err := hex.DecodeString("31303036343420524541444d452e6d640049866280e5a812f80c2750a43fe1ddc1f7225fc9343030303020626c6f6200488c2758caead757c828530bcfb0e0c8c223c0fb343030303020636d6400c1268e4d0283ae0e7c352a5e32cc824cfcf5aade31303036343420676f2e6d6f64002ae10c433ac16585d24ce94802e5e70839acb4da31303036343420676f2e73756d003712f87add624e7502d95b5224dd3ee89b3aca41343030303020696e64657800670c0ac1cf22b99ba436627a5a87517e5edf77c23430303030206c6f636b0014e606b6e3e773a4990ae6778e16f75858df8cf5313030363434206d61696e2e676f00e8a2072b578d9e710eb3fe2bd72d42b10a5fe6ab3430303030206f626a656374005153ad2bef7e18be52658983fdff82f5125163e834303030302072656600953fd4f4fc4ba2aa6680ae2dc9464c1359c17b103430303030207265706f7369746f7279000a897231e12c80dac24fcedde36d05a6262c0b0d3430303030207465737400f2ded005276610832df8850dd701c589b48feb913430303030207472656500a637ab746455126a45e7de80bade42e25df31f72")
	if err != nil {
		t.Fatalf("error creating test data from hex string: %v", err)
	}

	actual, err := DeserializeTree(data)
	if err != nil {
		t.Errorf("expected nil error but got: %v", err)
	}
	if actual.name != "" {
		t.Errorf("expected blank name but got: %s", actual.name)
	}
	if actual.oid != "" {
		t.Errorf("expected blank OID but got: %s", actual.oid)
	}

	if len(actual.entries) != 13 {
		t.Errorf("expected 13 entries but got: %d", len(actual.entries))
	}

	fmt.Println(actual)

	node, ok := actual.entries["main.go"]
	if !ok {
		t.Error("expected an entry for main.go but none exists")
	}
	if stubNode, ok := node.(stubNode); !ok {
		t.Errorf("expected an main.go to be a stub node but it is a %T", node)
	} else {
		if stubNode.name != "main.go" {
			t.Errorf("unexpected name: %s", stubNode.name)
		}
		if stubNode.oid != "e8a2072b578d9e710eb3fe2bd72d42b10a5fe6ab" {
			t.Errorf("unexpected oid: %s", stubNode.oid)
		}
		if stubNode.mode != "100644" {
			t.Errorf("unexpected mode: %s", stubNode.mode)
		}
	}

	node, ok = actual.entries["ref"]
	if !ok {
		t.Error("expected an entry for the ref dir but none exists")
	}
	if tree, ok := node.(*Tree); !ok {
		t.Errorf("expected an ref to be a tree but it is a %T", node)
	} else {
		if tree.name != "ref" {
			t.Errorf("unexpected name: %s", tree.name)
		}
		if tree.oid != "953fd4f4fc4ba2aa6680ae2dc9464c1359c17b10" {
			t.Errorf("unexpected oid: %s", tree.oid)
		}
		if tree.entries == nil {
			t.Errorf("tree entries should not be nil")
		}
		if len(tree.entries) > 0 {
			t.Errorf("tree entries should be empty")
		}
	}
}
