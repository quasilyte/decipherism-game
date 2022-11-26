package main

func checkAnagram(s1, s2 []byte) bool {
	if len(s1) != len(s2) {
		return false
	}
	var chars1 [256]byte
	var chars2 [256]byte
	for _, ch := range s1 {
		chars1[ch]++
	}
	for _, ch := range s2 {
		chars2[ch]++
	}
	for _, ch := range s1 {
		if chars1[ch] != chars2[ch] {
			return false
		}
	}
	return true
}

func rotateCharsRight(chars []byte) {
	if len(chars) == 0 {
		return
	}
	last := chars[len(chars)-1]
	for i := len(chars) - 1; i >= 1; i-- {
		chars[i] = chars[i-1]
	}
	chars[0] = last
}

func rotateCharsLeft(chars []byte) {
	if len(chars) == 0 {
		return
	}
	first := chars[0]
	for i := 0; i < len(chars)-1; i++ {
		chars[i] = chars[i+1]
	}
	chars[len(chars)-1] = first
}

func mapChars(chars []byte, f func(ch byte) byte) {
	for i, ch := range chars {
		chars[i] = f(ch)
	}
}

func mapEvenChars(chars []byte, f func(ch byte) byte) {
	for i, ch := range chars {
		pos := i + 1
		if pos%2 == 0 {
			chars[i] = f(ch)
		}
	}
}

func mapOddChars(chars []byte, f func(ch byte) byte) {
	for i, ch := range chars {
		pos := i + 1
		if pos%2 != 0 {
			chars[i] = f(ch)
		}
	}
}

func mapCharsButfirst(chars []byte, f func(ch byte) byte) {
	if len(chars) < 2 {
		return
	}
	data := chars[1:]
	for i, ch := range data {
		data[i] = f(ch)
	}
}

func mapCharsButlast(chars []byte, f func(ch byte) byte) {
	if len(chars) < 2 {
		return
	}
	data := chars[:len(chars)-1]
	for i, ch := range data {
		data[i] = f(ch)
	}
}

var dottedPairs = [256]byte{
	'a': 'z',
	'c': 'x',
	'e': 'v',
	'g': 't',
	'i': 'r',
	'k': 'p',
}

func polygraphicAtbash(chars []byte) {
	for i := 0; i < len(chars)-1; i++ {
		expectedNext := dottedPairs[chars[i]]
		if expectedNext == 0 {
			continue
		}
		if expectedNext == chars[i+1] {
			chars[i] = incChar(chars[i])
			chars[i+1] = decChar(chars[i+1])
			i++
		}
	}
}

var dottedChars = [256]bool{
	'a': true,
	'c': true,
	'e': true,
	'g': true,
	'i': true,
	'k': true,
	'p': true,
	'r': true,
	't': true,
	'v': true,
	'x': true,
	'z': true,
}

func incCharDotted(b byte) byte {
	if !dottedChars[b] {
		return b
	}
	return incChar(b)
}

func decCharUndotted(b byte) byte {
	if dottedChars[b] {
		return b
	}
	return decChar(b)
}

func incChar(b byte) byte {
	if b+1 > 'z' {
		return 'a'
	}
	return b + 1
}

func decChar(b byte) byte {
	if b-1 < 'a' {
		return 'z'
	}
	return b - 1
}

func (r *schemaRunner) decCharNowrap(b byte) byte {
	if b-1 < 'a' {
		return 'a'
	}
	return b - 1
}

func (r *schemaRunner) incCharNowrap(b byte) byte {
	if b+1 > 'z' {
		return 'z'
	}
	return b + 1
}
