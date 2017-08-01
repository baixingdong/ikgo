package ikgo

import (
	"unicode"
)

const (
	CHAR_USELESS   = 0
	CHAR_ARABIC    = 0x00000001
	CHAR_ENGLISH   = 0x00000002
	CHAR_CHINESE   = 0x00000004
	CHAR_OTHER_CJK = 0x00000008
)

/**
 * 识别字符类型
 * @param input
 * @return int CharacterUtil定义的字符类型常量
 */
func identifyCharType(input rune) int {
	if input > 47 && input < 58 {
		return CHAR_ARABIC
	}
	temp := input & (^0x20)
	if temp > 64 && temp < 91 {
		return CHAR_ENGLISH
	}
	if unicode.Is(unicode.Scripts["Han"], input) {
		return CHAR_CHINESE
	}
	if unicode.Is(unicode.Scripts["Hangul"], input) || unicode.Is(unicode.Scripts["Hiragana"], input) || unicode.Is(unicode.Scripts["Katakana"], input) {
		return CHAR_OTHER_CJK
	}

	return CHAR_USELESS
}

func regularize(input byte, lowercase bool) (output byte) {
	// 全半角以及大小写转换会放到前面，这里不做处理
	output = input
	return
}
