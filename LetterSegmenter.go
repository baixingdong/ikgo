package ikgo

var (
	Letter_Connector map[rune]bool
	Num_Connector    map[rune]bool
)

type LetterSegmenter struct {
	name                     string
	start, end               int
	englishStart, englishEnd int
	arabicStart, arabicEnd   int
}

func initLS() {
	ca := []byte{'#', '&', '+', '-', '.', '@', '_'}
	cb := []byte{',', '.'}
	Letter_Connector = make(map[rune]bool)
	Num_Connector = make(map[rune]bool)
	for _, a := range ca {
		Letter_Connector[rune(a)] = true
	}

	for _, b := range cb {
		Num_Connector[rune(b)] = true
	}
}

func NewLetterSegmenter() *LetterSegmenter {
	return &LetterSegmenter{
		start:        -1,
		end:          -1,
		englishStart: -1,
		englishEnd:   -1,
		arabicStart:  -1,
		arabicEnd:    -1,
	}
}

/**
 * 判断是否是字母连接符号
 * @param input
 * @return
 */
func (s *LetterSegmenter) isLetterConnector(input rune) bool {
	_, exists := Letter_Connector[input]
	return exists
}

/**
 * 判断是否是数字连接符号
 * @param input
 * @return
 */
func (s *LetterSegmenter) isNumConnector(input rune) bool {
	_, exists := Num_Connector[input]
	return exists
}

/**
 * 处理数字字母混合输出
 * 如：windos2000 | linliangyi2005@gmail.com
 * @param input
 * @param context
 * @return
 */
func (s *LetterSegmenter) processMixLetter(context *AnalyzeContext) bool {
	needLock := false

	if s.start == -1 { //当前的分词器尚未开始处理字符
		if CHAR_ARABIC == context.charType[context.cursor] || CHAR_ENGLISH == context.charType[context.cursor] {
			//记录起始指针的位置,标明分词器进入处理状态
			s.start = context.cursor
			s.end = s.start
		}
	} else { //当前的分词器正在处理字符
		if CHAR_ARABIC == context.charType[context.cursor] || CHAR_ENGLISH == context.charType[context.cursor] {
			//记录下可能的结束位置
			s.end = context.cursor
		} else if CHAR_USELESS == context.charType[context.cursor] && s.isLetterConnector(context.segmentBuff[context.cursor]) {
			//记录下可能的结束位置
			s.end = context.cursor
		} else {
			//遇到非Letter字符，输出词元
			newLexeme := NewLexeme(context.bufOffset, s.start, s.end-s.start+1, LEXEME_TYPE_LETTER)
			context.addLexeme(newLexeme)
			s.start = -1
			s.end = -1
		}
	}

	//判断缓冲区是否已经读完
	if context.isBufferConsumed() && s.start != -1 && s.end != -1 {
		//缓冲以读完，输出词元
		newLexeme := NewLexeme(context.bufOffset, s.start, s.end-s.start+1, LEXEME_TYPE_LETTER)
		context.addLexeme(newLexeme)
		s.start = -1
		s.end = -1
	}

	//判断是否锁定缓冲区
	if s.start == -1 && s.end == -1 {
		//对缓冲区解锁
		needLock = false
	} else {
		needLock = true
	}
	return needLock
}

/**
 * 处理纯英文字母输出
 * @param context
 * @return
 */
func (s *LetterSegmenter) processEnglishLetter(context *AnalyzeContext) bool {
	needLock := false
	if s.englishStart == -1 { //当前的分词器尚未开始处理英文字符
		if CHAR_ENGLISH == context.charType[context.cursor] {
			//记录起始指针的位置,标明分词器进入处理状态
			s.englishStart = context.cursor
			s.englishEnd = s.englishStart
		}
	} else { //当前的分词器正在处理英文字符
		if CHAR_ENGLISH == context.charType[context.cursor] {
			//记录当前指针位置为结束位置
			s.englishEnd = context.cursor
		} else {
			//遇到非English字符,输出词元
			newLexeme := NewLexeme(context.bufOffset, s.englishStart, s.englishEnd-s.englishStart+1, LEXEME_TYPE_ENGLISH)
			context.addLexeme(newLexeme)
			s.englishStart = -1
			s.englishEnd = -1
		}
	}

	//判断缓冲区是否已经读完
	if context.isBufferConsumed() && s.englishStart != -1 && s.englishEnd != -1 {
		//缓冲已读完，输出词元
		newLexeme := NewLexeme(context.bufOffset, s.englishStart, s.englishEnd-s.englishStart+1, LEXEME_TYPE_ENGLISH)
		context.addLexeme(newLexeme)
		s.start = -1
		s.end = -1
	}

	//判断是否锁定缓冲区
	if s.englishStart == -1 && s.englishEnd == -1 {
		//对缓冲区解锁
		needLock = false
	} else {
		needLock = true
	}
	return needLock
}

/**
 * 处理阿拉伯数字输出
 * @param context
 * @return
 */
func (s *LetterSegmenter) processArabicLetter(context *AnalyzeContext) bool {
	needLock := false

	if s.arabicStart == -1 { //当前的分词器尚未开始处理数字字符
		if CHAR_ARABIC == context.charType[context.cursor] {
			//记录起始指针的位置,标明分词器进入处理状态
			s.arabicStart = context.cursor
			s.arabicEnd = s.arabicStart
		}
	} else {
		//当前的分词器正在处理数字字符
		if CHAR_ARABIC == context.charType[context.cursor] {
			//记录当前指针位置为结束位置
			s.arabicEnd = context.cursor
		} else if CHAR_USELESS == context.charType[context.cursor] && s.isNumConnector(context.segmentBuff[context.cursor]) {
			//不输出数字，但不标记结束
		} else {
			//遇到非Arabic字符,输出词元
			newLexeme := NewLexeme(context.bufOffset, s.arabicStart, s.arabicEnd-s.arabicStart+1, LEXEME_TYPE_ARABIC)
			context.addLexeme(newLexeme)
			s.arabicStart = -1
			s.arabicEnd = -1
		}
	}

	//判断缓冲区是否已经读完
	if context.isBufferConsumed() && s.arabicStart != -1 && s.arabicEnd != -1 {
		//缓冲以读完，输出词元
		newLexeme := NewLexeme(context.bufOffset, s.arabicStart, s.arabicEnd-s.arabicStart+1, LEXEME_TYPE_ARABIC)
		context.addLexeme(newLexeme)
		s.start = -1
		s.end = -1
	}

	//判断是否锁定缓冲区
	if s.arabicStart == -1 && s.arabicEnd == -1 {
		//对缓冲区解锁
		needLock = false
	} else {
		needLock = true
	}
	return needLock
}

func (s *LetterSegmenter) analyze(context *AnalyzeContext) {
	bufferLockFlag := false
	//处理英文字母
	bufferLockFlag = s.processEnglishLetter(context) || bufferLockFlag
	//处理阿拉伯字母
	bufferLockFlag = s.processArabicLetter(context) || bufferLockFlag
	//处理混合字母(这个要放最后处理，可以通过QuickSortSet排除重复)
	bufferLockFlag = s.processMixLetter(context) || bufferLockFlag

	//判断是否锁定缓冲区
	if bufferLockFlag {
		context.lockBuffer(s.name)
	} else {
		//对缓冲区解锁
		context.unlockBuffer(s.name)
	}
}

func (s *LetterSegmenter) reset() {
	s.start = -1
	s.end = -1
	s.englishStart = -1
	s.englishEnd = -1
	s.arabicStart = -1
	s.arabicEnd = -1
}
