package say

import (
	"fmt"
	"strings"
)

func checkOut(out string, stream bool) error {
	if out == "" {
		return nil
	}
	n, bad := countIntVerbs(out)
	if stream {
		if bad || n != 1 {
			return fmt.Errorf("say: -o needs one %%d in --stream")
		}
		return nil
	}
	if bad || n > 0 {
		return fmt.Errorf("say: -o template needs --stream")
	}
	return nil
}

func outPath(out string, idx int) string {
	if out == "" {
		return ""
	}
	n, bad := countIntVerbs(out)
	if !bad && n == 1 {
		return fmt.Sprintf(out, idx)
	}
	return strings.ReplaceAll(out, "%%", "%")
}

func countIntVerbs(s string) (n int, bad bool) {
	r := []rune(s)
	for i := 0; i < len(r); i++ {
		if r[i] != '%' {
			continue
		}
		if i+1 >= len(r) {
			return n, true
		}
		i++
		if r[i] == '%' {
			continue
		}
		if r[i] == '-' {
			i++
			if i >= len(r) {
				return n, true
			}
		}
		for i < len(r) && r[i] >= '0' && r[i] <= '9' {
			i++
		}
		if i >= len(r) || r[i] != 'd' {
			return n, true
		}
		n++
	}
	return n, false
}
