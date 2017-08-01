package ikgo

/**
 * IK词元对象
 */

const (
	LEXEME_TYPE_UNKNOWN   = 0
	LEXEME_TYPE_ENGLISH   = 1
	LEXEME_TYPE_ARABIC    = 2
	LEXEME_TYPE_LETTER    = 3
	LEXEME_TYPE_CNWORD    = 4
	LEXEME_TYPE_CNCHAR    = 64
	LEXEME_TYPE_OTHER_CJK = 8
	LEXEME_TYPE_CNUM      = 16
	LEXEME_TYPE_COUNT     = 32
	LEXEME_TYPE_CQUAN     = 48
)

type Lexeme struct {
	offset, begin, length int
	lexemeText            string
	lexemeType            int
}

func NewLexeme(offset, begin, length, lexemeType int) (l *Lexeme) {
	if length < 0 {
		l = nil
		return
	}
	l = &Lexeme{
		offset:     offset,
		begin:      begin,
		length:     length,
		lexemeType: lexemeType,
	}
	return
}

/*
 * 判断词元相等算法
 * 起始位置偏移、起始位置、终止位置相同
 */
func (l *Lexeme) equals(o *Lexeme) (b bool) {
	if o == nil {
		b = false
		return
	}
	if o == l {
		b = true
		return
	}
	if l.offset == o.offset && l.begin == o.begin && l.length == o.length {
		b = true
		return
	}
	b = false
	return
}

/*
 * 词元哈希编码算法
 */
func (l *Lexeme) hash() int {
	beg := l.begin + l.offset
	end := beg + l.length
	return (beg * 37) + (end * 31) + ((beg*end)%l.length)*11
}

/*
 * 词元在排序集合中的比较算法
 */

func (l *Lexeme) compare(o *Lexeme) int {
	if l.begin < o.begin {
		return -1
	}
	if l.begin > o.begin {
		return 1
	}
	if l.length > o.length {
		return -1
	}
	if l.length < o.length {
		return 1
	}
	return 0
}

/**
 * 获取词元在文本中的起始位置
 * @return int
 */
func (l *Lexeme) GetBeginPosition() int {
	return l.offset + l.begin
}

/**
 * 获取词元在文本中的结束位置
 * @return int
 */
func (l *Lexeme) GetEndPosition() int {
	return l.offset + l.begin + l.length
}

/**
 * 获取词元长度
 * @return int
 */
func (l *Lexeme) GetLength() int {
	return l.length
}

func (l *Lexeme) GetTypeString() string {
	switch l.lexemeType {
	case LEXEME_TYPE_ENGLISH:
		return "ENGLISH"
	case LEXEME_TYPE_ARABIC:
		return "ARABIC"
	case LEXEME_TYPE_LETTER:
		return "LETTER"
	case LEXEME_TYPE_CNWORD:
		return "CN_WORD"
	case LEXEME_TYPE_CNCHAR:
		return "CN_CHAR"
	case LEXEME_TYPE_OTHER_CJK:
		return "OTHER_CJK"
	case LEXEME_TYPE_COUNT:
		return "COUNT"
	case LEXEME_TYPE_CNUM:
		return "TYPE_CNUM"
	case LEXEME_TYPE_CQUAN:
		return "TYPE_CQUAN"
	default:
		return "UNKNOWN"
	}
}

/**
 * 合并两个相邻的词元
 * @param l
 * @param lexemeType
 * @return boolean 词元是否成功合并
 */

func (l *Lexeme) append(o *Lexeme, lexemeType int) bool {
	if o == nil {
		return false
	}
	if l.GetEndPosition() != o.GetBeginPosition() {
		return false
	}
	l.length += o.length
	l.lexemeType = lexemeType
	return true
}
