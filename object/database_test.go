package object

import (
	"compress/zlib"
	"io/ioutil"
	"os"
	"testing"
)

func setUpTestDatabase(t *testing.T) Database {
	dbDir, err := ioutil.TempDir("", "got_test_db_*")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}

	return NewDatabase(dbDir)
}

func TestStoreAndRetreiveGenericStorable(t *testing.T) {
	db := setUpTestDatabase(t)

	obj := &genericStorable{
		storableType: "blob",
		size:         13,
		data:         []byte("hello, world!"),
	}

	actualOID, err := db.Store(obj)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	expectedOID := "1f09d30c707d53f3d16c530dd73d70a6ce7596a9"
	if actualOID != expectedOID {
		t.Errorf("expected OID '%s' but got '%s'", "abc123", actualOID)
	}

	f, err := os.Open(db.(*database).objectPath(actualOID))
	if err != nil {
		t.Fatalf("unable to open data file: %v", err)
	}

	unzipper, err := zlib.NewReader(f)
	if err != nil {
		t.Fatalf("unable to open zlib reader: %v", err)
	}

	actualData, err := ioutil.ReadAll(unzipper)
	if err != nil {
		t.Fatalf("expected data file '%s' to exist, but could not read: %v", db.(*database).objectPath(actualOID), err)
	}
	expectedData := "blob 13\x00hello, world!"
	if string(actualData) != expectedData {
		t.Errorf("unexpected data stored on disk: %v", actualData)
	}

	readObj, err := db.Read(actualOID)
	if err != nil {
		t.Fatalf("error reading back object from database: %v", err)
	}

	if readObj.Type() != "blob" {
		t.Errorf("expected type 'blob' but got '%s'", readObj.Type())
	}

	dataString := string(readObj.Serialize())
	if dataString != "hello, world!" {
		t.Errorf("expected data 'hello, world!' but got '%s'", string(readObj.Serialize()))
	}
}
