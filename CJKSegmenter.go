package ikgo

import (
	"container/list"
)

type CJKSegmenter struct {
	name    string
	tmpHits *list.List
}

func NewCJKSegmenter() *CJKSegmenter {
	return &CJKSegmenter{name: "CJK_SEGMENTER", tmpHits: list.New()}
}

func (s *CJKSegmenter) analyze(context *AnalyzeContext) {
	if CHAR_USELESS != context.charType[context.cursor] {
		//优先处理tmpHits中的hit
		if s.tmpHits.Len() != 0 {
			//处理词段队列
			for iter := s.tmpHits.Front(); iter != nil; {
				cur := iter
				iter = iter.Next()
				hit := cur.Value.(*Hit)
				hit = matchWithHit(context.segmentBuff, context.cursor, hit)
				if hit.isMatch() {
					//输出当前的词
					newLexeme := NewLexeme(context.bufOffset, hit.beg, context.cursor-hit.beg+1, LEXEME_TYPE_CNWORD)
					context.addLexeme(newLexeme)

					if !hit.isPrefix() { //不是词前缀，hit不需要继续匹配，移除
						s.tmpHits.Remove(cur)
					}
				} else if hit.isUnmatch() {
					//hit不是词，移除
					s.tmpHits.Remove(cur)
				}

			}
		}
		//*********************************
		//再对当前指针位置的字符进行单字匹配
		singleCharHit := MainDict.matchSeg(context.segmentBuff, context.cursor, 1)
		if singleCharHit.isMatch() { //首字成词
			//输出当前的词
			newLexeme := NewLexeme(context.bufOffset, context.cursor, 1, LEXEME_TYPE_CNWORD)
			context.addLexeme(newLexeme)
			//同时也是词前缀
			if singleCharHit.isPrefix() {
				//前缀匹配则放入hit列表
				s.tmpHits.PushBack(singleCharHit)
			}
		} else if singleCharHit.isPrefix() { //首字为词前缀
			//前缀匹配则放入hit列表
			s.tmpHits.PushBack(singleCharHit)
		}
	} else {
		//遇到CHAR_USELESS字符
		//清空队列
		s.tmpHits = list.New()
	}

	//判断缓冲区是否已经读完
	if context.isBufferConsumed() {
		//清空队列
		s.tmpHits = list.New()
	}

	//判断是否锁定缓冲区
	if s.tmpHits.Len() == 0 {
		context.unlockBuffer(s.name)
	} else {
		context.lockBuffer(s.name)
	}

}

func (s *CJKSegmenter) reset() {
	s.tmpHits = list.New()
}
