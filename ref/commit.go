package ref

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TypeCommit is the type returned by Commit objects.
const TypeCommit = "commit"

var authorRegexp = regexp.MustCompile(`^(.*) <(.*)> ([0-9]+) -([0-9]{2})([0-9]{2})$`)

type Author struct {
	Name  string
	Email string
	Time  time.Time
}

func (a Author) String() string {
	return fmt.Sprintf("%s <%s> %d %s", a.Name, a.Email, a.Time.Unix(), a.Time.Format("-0700"))
}

type Commit struct {
	Parent  string
	TreeOID string
	Author  Author
	Message string
}

func NewCommit(parent string, treeOID string, author Author, message string) Commit {
	return Commit{
		Parent:  parent,
		TreeOID: treeOID,
		Author:  author,
		Message: message,
	}
}

func DeserializeCommit(data []byte) (result Commit, err error) {
	r := bufio.NewReader(bytes.NewBuffer(data))

	var line string
	var headers = map[string]string{}
	for err == nil {
		line, err = r.ReadString('\n')
		if line == "\n" {
			break
		}

		split := strings.Index(line, " ")
		if split == -1 {
			err = fmt.Errorf("invalid commit format: '%s'", line)
			break
		}

		headers[line[:split]] = strings.TrimSpace(line[split+1:])
	}
	if err != nil {
		return result, fmt.Errorf("error scanning commit data: %w", err)
	}

	var message string
	for err == nil {
		line, err = r.ReadString('\n')
		message += line
	}
	if err != nil && err != io.EOF {
		return result, fmt.Errorf("error parsing commit message: %w", err)
	}
	err = nil

	author, err := parseAuthorString(headers["author"])
	if err != nil {
		return result, err
	}

	return NewCommit(headers["parent"], headers["tree"], author, message), nil
}

func (c Commit) Type() string {
	return "commit"
}

func (c Commit) Serialize() []byte {
	var parentLine string
	if c.Parent != "" {
		parentLine = fmt.Sprintf("parent %s\n", c.Parent)
	}

	data := fmt.Sprintf("tree %s\n%sauthor %s\ncommitter %s\n\n%s\n", c.TreeOID, parentLine, c.Author.String(), c.Author.String(), c.Message)

	return []byte(data)
}

func parseAuthorString(input string) (result Author, err error) {
	m := authorRegexp.FindStringSubmatch(input)
	if len(m) != 6 {
		return result, fmt.Errorf("invalid author format: '%s'", input)
	}

	result.Name = m[1]
	result.Email = m[2]

	tstamp, _ := strconv.ParseInt(m[3], 10, 64)
	hoursOffset, _ := strconv.Atoi(m[4])
	minsOffset, _ := strconv.Atoi(m[4])

	result.Time = time.Unix(tstamp, 0).In(time.FixedZone("", -(hoursOffset*60*60 + minsOffset*60)))
	return
}
