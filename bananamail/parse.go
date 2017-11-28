package main

func ParseMsg(msg string) string {
	var s string
	for _, a := range msg {
		if a == 45 {
			s += "\n"
			continue
		}
		s += string(a)
	}
	return s
}
