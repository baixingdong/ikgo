package ikgo

import "container/list"

var (
	Chn_Num        = []rune("一二两三四五六七八九十零壹贰叁肆伍陆柒捌玖拾百千万亿拾佰仟萬億兆卅廿")
	ChnNumberChars map[rune]bool
)

type CN_QuantifierSegmenter struct {
	name string
	/*
	 * 词元的开始位置，
	 * 同时作为子分词器状态标识
	 * 当start > -1 时，标识当前的分词器正在处理字符
	 */
	nStart int
	/*
	 * 记录词元结束位置
	 * end记录的是在词元中最后一个出现的合理的数词结束
	 */
	nEnd      int
	countHits *list.List
}

func initCNQS() {
	ChnNumberChars = make(map[rune]bool)
	for _, num := range Chn_Num {
		ChnNumberChars[num] = true
	}
}

func NewCN_QuantifierSegmenter() *CN_QuantifierSegmenter {
	return &CN_QuantifierSegmenter{nStart: -1, nEnd: -1, countHits: list.New(), name: "QUAN_SEGMENTER"}
}

/**
 * 添加数词词元到结果集
 * @param context
 */
func (s *CN_QuantifierSegmenter) outputNumLexeme(context *AnalyzeContext) {
	if s.nStart > -1 && s.nEnd > -1 {
		//输出数词
		newLexeme := NewLexeme(context.bufOffset, s.nStart, s.nEnd-s.nStart+1, LEXEME_TYPE_CNUM)
		context.addLexeme(newLexeme)
	}
}

/**
 * 处理数词
 */
func (s *CN_QuantifierSegmenter) processCNumber(context *AnalyzeContext) {
	if s.nStart == -1 && s.nEnd == -1 { //初始状态
		if CHAR_CHINESE == context.charType[context.cursor] {
			if _, exists := ChnNumberChars[context.segmentBuff[context.cursor]]; exists {
				//记录数词的起始、结束位置
				s.nStart = context.cursor
				s.nEnd = context.cursor
			}
		}
	} else { //正在处理状态
		if CHAR_CHINESE == context.charType[context.cursor] {
			if _, exists := ChnNumberChars[context.segmentBuff[context.cursor]]; exists {
				//记录数词的结束位置
				s.nEnd = context.cursor
			}
		} else {
			//输出数词
			s.outputNumLexeme(context)
			//重置头尾指针
			s.nStart = -1
			s.nEnd = -1
		}
	}

	//缓冲区已经用完，还有尚未输出的数词
	if context.isBufferConsumed() && s.nStart != -1 && s.nEnd != -1 {
		//输出数词
		s.outputNumLexeme(context)
		//重置头尾指针
		s.nStart = -1
		s.nEnd = -1
	}
}

/**
 * 判断是否需要扫描量词
 * @return
 */
func (s *CN_QuantifierSegmenter) needCountScan(context *AnalyzeContext) bool {
	if (s.nStart != -1 && s.nEnd != -1) || s.countHits.Len() != 0 {
		//正在处理中文数词,或者正在处理量词
		return true
	}
	//找到一个相邻的数词
	if context.getOrgLexemes().size != 0 {
		l := context.getOrgLexemes().peekLast()
		if (LEXEME_TYPE_CNUM == l.lexemeType || LEXEME_TYPE_ARABIC == l.lexemeType) && (l.begin+l.length == context.cursor) {
			return true
		}
	}
	return false
}

/**
 * 处理中文量词
 * @param context
 */
func (s *CN_QuantifierSegmenter) processCount(context *AnalyzeContext) {
	// 判断是否需要启动量词扫描
	if !s.needCountScan(context) {
		return
	}

	if CHAR_CHINESE == context.charType[context.cursor] {
		//优先处理countHits中的hit
		if s.countHits.Len() != 0 {
			//处理词段队列
			for iter := s.countHits.Front(); iter != nil; {
				cur := iter
				iter = iter.Next()
				hit := cur.Value.(*Hit)
				hit = matchWithHit(context.segmentBuff, context.cursor, hit)
				if hit.isMatch() {
					//输出当前的词
					newLexeme := NewLexeme(context.bufOffset, hit.beg, context.cursor-hit.beg+1, LEXEME_TYPE_COUNT)
					context.addLexeme(newLexeme)

					if !hit.isPrefix() { //不是词前缀，hit不需要继续匹配，移除
						s.countHits.Remove(cur)
					}
				} else if hit.isUnmatch() {
					//hit不是词，移除
					s.countHits.Remove(cur)
				}
			}
		}

		//*********************************
		//再对当前指针位置的字符进行单字匹配
		singleCharHit := MainDict.matchSeg(context.segmentBuff, context.cursor, 1)
		if singleCharHit.isMatch() { //首字成量词词
			//输出当前的词
			newLexeme := NewLexeme(context.bufOffset, context.cursor, 1, LEXEME_TYPE_COUNT)
			context.addLexeme(newLexeme)
			//同时也是词前缀
			if singleCharHit.isPrefix() {
				//前缀匹配则放入hit列表
				s.countHits.PushBack(singleCharHit)
			}
		} else if singleCharHit.isPrefix() { //首字为词前缀
			//前缀匹配则放入hit列表
			s.countHits.PushBack(singleCharHit)
		}
	} else {
		//输入的不是中文字符
		//清空未成形的量词
		s.countHits = list.New()
	}
	//判断缓冲区是否已经读完
	if context.isBufferConsumed() {
		//清空队列
		s.countHits = list.New()
	}
}

/**
 * 分词
 */
func (s *CN_QuantifierSegmenter) analyze(context *AnalyzeContext) {
	//处理中文数词
	s.processCNumber(context)
	//处理中文量词
	s.processCount(context)
	//判断是否锁定缓冲区
	if s.nStart == -1 && s.nEnd == -1 && s.countHits.Len() == 0 {
		//对缓冲区解锁
		context.unlockBuffer(s.name)
	} else {
		context.lockBuffer(s.name)
	}
}

func (s *CN_QuantifierSegmenter) reset() {
	s.nStart = -1
	s.nEnd = -1
	s.countHits = list.New()
}
