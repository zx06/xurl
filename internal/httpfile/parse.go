package httpfile

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/zx06/xurl/internal/errors"
)

type Request struct {
	Name    string
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Vars    map[string]string // @name=value from file
}

func Parse(content []byte) ([]Request, *errors.XError) {
	blocks := splitBlocks(content)
	if len(blocks) == 0 {
		return nil, errors.New(errors.CodeHTTPFile, "no HTTP requests found in file", nil)
	}
	out := make([]Request, 0, len(blocks))
	for _, block := range blocks {
		req, xe := parseBlock(block.text, block.name)
		if xe != nil {
			return nil, xe
		}
		out = append(out, req)
	}
	return out, nil
}

type block struct {
	name string
	text []byte
}

func splitBlocks(content []byte) []block {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var blocks []block
	var curName string
	var curLines []string

	flush := func() {
		if len(curLines) == 0 {
			return
		}
		text := strings.Join(curLines, "\n")
		blocks = append(blocks, block{name: curName, text: []byte(text)})
		curLines = nil
		curName = ""
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "###") {
			flush()
			name := strings.TrimSpace(strings.TrimPrefix(line, "###"))
			curName = name
			continue
		}
		curLines = append(curLines, line)
	}
	flush()
	return blocks
}

func parseBlock(text []byte, name string) (Request, *errors.XError) {
	lines := strings.Split(string(text), "\n")
	req := Request{
		Name:    name,
		Headers: map[string]string{},
		Vars:    map[string]string{},
	}

	i := 0
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}
	for i < len(lines) {
		trim := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trim, "@") && strings.Contains(trim, "=") {
			eq := strings.Index(trim, "=")
			key := strings.TrimPrefix(strings.TrimSpace(trim[:eq]), "@")
			val := strings.TrimSpace(trim[eq+1:])
			if key != "" {
				req.Vars[key] = val
			}
			i++
			continue
		}
		break
	}

	if i >= len(lines) {
		return Request{}, errors.New(errors.CodeHTTPFile, "missing request line", map[string]any{"name": name})
	}

	reqLine := strings.TrimSpace(lines[i])
	i++
	method, url, xe := parseRequestLine(reqLine)
	if xe != nil {
		return Request{}, xe
	}
	req.Method = method
	req.URL = url

	var bodyLines []string
	inBody := false
	for i < len(lines) {
		line := lines[i]
		if !inBody {
			if strings.TrimSpace(line) == "" {
				inBody = true
				i++
				continue
			}
			if hName, hVal, ok := parseHeaderLine(line); ok {
				req.Headers[hName] = hVal
				i++
				continue
			}
			// Non-empty non-header before blank line — treat as body start
			inBody = true
			bodyLines = append(bodyLines, line)
			i++
			continue
		}
		bodyLines = append(bodyLines, line)
		i++
	}
	if len(bodyLines) > 0 {
		req.Body = strings.Join(bodyLines, "\n")
		if strings.HasSuffix(req.Body, "\n") {
			// keep as-is from join without trailing newline unless file had it
		}
	}
	return req, nil
}

func parseRequestLine(line string) (method, url string, xe *errors.XError) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return "", "", errors.New(errors.CodeHTTPFile, "invalid request line", map[string]any{"line": line})
	}
	method = strings.ToUpper(parts[0])
	url = parts[1]
	if len(parts) >= 3 && strings.HasPrefix(parts[2], "HTTP/") {
		// drop HTTP version
	}
	return method, url, nil
}

func parseHeaderLine(line string) (name, value string, ok bool) {
	idx := strings.Index(line, ":")
	if idx <= 0 {
		return "", "", false
	}
	name = strings.TrimSpace(line[:idx])
	value = strings.TrimSpace(line[idx+1:])
	if name == "" {
		return "", "", false
	}
	return name, value, true
}

func FindByName(requests []Request, name string) (Request, bool) {
	if name == "" {
		if len(requests) == 1 {
			return requests[0], true
		}
		return Request{}, false
	}
	for _, r := range requests {
		if r.Name == name {
			return r, true
		}
	}
	return Request{}, false
}
