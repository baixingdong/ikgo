package ikgo

import (
	"fmt"
	"strings"
)

/**
 * Lexeme链（路径）
 */
type LexemePath struct {
	set                               QuickSortSet
	pathBegin, pathEnd, payloadLength int // 起止位置以及词元链的有效字符长度
}

func NewLexemePath() (l *LexemePath) {
	l = &LexemePath{pathBegin: -1, pathEnd: -1, payloadLength: 0}
	return
}

/**
 * 检测词元位置交叉（有歧义的切分）
 * @param lexeme
 * @return
 */
func (lp *LexemePath) checkCross(l *Lexeme) bool {
	return (l.begin >= lp.pathBegin && l.begin < lp.pathEnd) ||
		(lp.pathBegin >= l.begin && lp.pathBegin < l.begin+l.length)

}

/**
 * 向LexemePath追加相交的Lexeme
 * @param lexeme
 * @return
 */
func (lp *LexemePath) addCrossLexeme(l *Lexeme) bool {
	if lp.set.size == 0 {
		lp.set.addLexeme(l)
		lp.pathBegin = l.begin
		lp.pathEnd = l.begin + l.length
		lp.payloadLength += l.length
		return true
	}
	if lp.checkCross(l) {
		lp.set.addLexeme(l)
		if l.begin+l.length > lp.pathEnd {
			lp.pathEnd = l.begin + l.length
		}
		lp.payloadLength = lp.pathEnd - lp.pathBegin
		return true
	}
	return false
}

/**
 * 向LexemePath追加不相交的Lexeme
 * @param lexeme
 * @return
 */
func (lp *LexemePath) addNotCrossLexeme(l *Lexeme) bool {
	if lp.set.size == 0 {
		lp.set.addLexeme(l)
		lp.pathBegin = l.begin
		lp.pathEnd = l.begin + l.length
		lp.payloadLength += l.length
		return true
	}
	if lp.checkCross(l) {
		return false
	}
	lp.set.addLexeme(l)
	lp.payloadLength += l.length
	head := lp.set.peekFirst()
	tail := lp.set.peekLast()
	lp.pathBegin = head.begin
	lp.pathEnd = tail.begin + tail.length
	return true
}

/**
 * 移除尾部的Lexeme
 * @return
 */
func (lp *LexemePath) removeTail() (l *Lexeme) {
	l = lp.set.pollLast()
	if lp.set.size == 0 {
		lp.pathBegin = -1
		lp.pathEnd = -1
		lp.payloadLength = 0
	} else {
		lp.payloadLength -= l.length
		tail := lp.set.peekLast()
		lp.pathEnd = tail.begin + tail.length
	}

	return
}

/**
 * 获取LexemePath的路径长度
 * @return
 */
func (lp *LexemePath) getPathLength() int {
	return lp.pathEnd - lp.pathBegin
}

/**
 * X权重（词元长度积）
 * @return
 */

func (lp *LexemePath) getXWeight() (product int) {
	product = 1
	c := lp.set.head
	for c != nil && c.lexeme != nil {
		product *= c.lexeme.length
		c = c.next
	}
	return
}

/**
 * 词元位置权重
 * @return
 */
func (lp *LexemePath) getPWeight() (product int) {
	product = 0
	p := 0
	c := lp.set.head
	for c != nil && c.lexeme != nil {
		p++
		product += p * c.lexeme.length
		c = c.next
	}
	return
}

func (lp *LexemePath) deepCopy() *LexemePath {
	nlp := &LexemePath{}
	nlp.pathBegin = lp.pathBegin
	nlp.pathEnd = lp.pathEnd
	nlp.payloadLength = lp.payloadLength
	c := lp.set.head
	for c != nil && c.lexeme != nil {
		nlp.set.addLexeme(c.lexeme)
		c = c.next
	}

	return nlp
}

func (lp *LexemePath) compare(nlp *LexemePath) int {
	if lp.payloadLength > nlp.payloadLength {
		return -1
	}
	if lp.payloadLength < nlp.payloadLength {
		return 1
	}
	if lp.set.size < nlp.set.size {
		return -1
	}
	if lp.set.size > nlp.set.size {
		return 1
	}
	if lp.getPathLength() > nlp.getPathLength() {
		return -1
	}
	if lp.getPathLength() < nlp.getPathLength() {
		return 1
	}
	if lp.pathEnd > nlp.pathEnd {
		return -1
	}
	if lp.pathEnd < nlp.pathEnd {
		return 1
	}
	xw1 := lp.getXWeight()
	xw2 := nlp.getXWeight()
	if xw1 > xw2 {
		return -1
	}
	if xw1 < xw2 {
		return 1
	}
	pw1 := lp.getPWeight()
	pw2 := nlp.getPWeight()
	if pw1 > pw2 {
		return -1
	}
	if pw1 < pw2 {
		return 1
	}
	return 0
}

func (lp *LexemePath) toString() string {
	sb := fmt.Sprintf("pathBegin: %d", lp.pathBegin)
	se := fmt.Sprintf("pathEnd : %d", lp.pathEnd)
	sp := fmt.Sprintf("pathPayload: %d", lp.payloadLength)
	sl := []string{sb, se, sp}
	head := lp.set.head
	for head != nil {
		sl = append(sl, fmt.Sprintf("lexme: %+v", head.lexeme))
		head = head.next
	}
	return strings.Join(sl, "\n")
}
