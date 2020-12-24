package ref

import (
	"fmt"
	"time"
)

type Author struct {
	Name  string
	Email string
	Time  time.Time
}

func (a Author) String() string {
	return fmt.Sprintf("%s <%s> %d %s", a.Name, a.Email, a.Time.Unix(), a.Time.Format("-0700"))
}

type Commit struct {
	parent  string
	oid     string
	author  Author
	message string
}

func NewCommit(parent string, oid string, author Author, message string) Commit {
	return Commit{
		parent:  parent,
		oid:     oid,
		author:  author,
		message: message,
	}
}

func (c Commit) Type() string {
	return "commit"
}

func (c Commit) Serialize() []byte {
	var parentLine string
	if c.parent != "" {
		parentLine = fmt.Sprintf("parent %s\n", c.parent)
	}

	data := fmt.Sprintf("tree %s\n%sauthor %s\ncommitter %s\n\n%s\n", c.oid, parentLine, c.author.String(), c.author.String(), c.message)

	return append([]byte(fmt.Sprintf("%s %d\x00", c.Type(), len(data))), []byte(data)...)
}
