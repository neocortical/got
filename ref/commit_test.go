package ref

import "testing"

func TestDeserializeCommit(t *testing.T) {
	data := []byte(`tree 0e3d6d78ab2bce1cfdcdc9c4f745f186c8b6daa7
parent bccd3e06dd549a5c27497f6a11243019ba2abb80
author Nathan Smith <nathan@neocortical.net> 1609095922 -0800
committer Nathan Smith <nathan@neocortical.net> 1609095922 -0800

commit message
on
three lines`)

	actual, err := DeserializeCommit(data)
	if err != nil {
		t.Errorf("expected nil error but got: %v", err)
	}
	if actual.TreeOID != "0e3d6d78ab2bce1cfdcdc9c4f745f186c8b6daa7" {
		t.Errorf("unexpected value for tree OID: %s", actual.TreeOID)
	}
	if actual.Parent != "bccd3e06dd549a5c27497f6a11243019ba2abb80" {
		t.Errorf("unexpected value for parent OID: %s", actual.Parent)
	}
	if actual.Author.Name != "Nathan Smith" {
		t.Errorf("unexpected value for author name: %s", actual.Author.Name)
	}
	if actual.Author.Email != "nathan@neocortical.net" {
		t.Errorf("unexpected value for author email: %s", actual.Author.Name)
	}
	if actual.Author.Time.Unix() != 1609095922 {
		t.Errorf("unexpected value for author commit time: %d", actual.Author.Time.Unix())
	}
	if actual.Message != "commit message\non\nthree lines" {
		t.Errorf("unexpected value for commit message: %s", actual.Message)
	}
}
