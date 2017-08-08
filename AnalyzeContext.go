package ikgo

import (
	"bufio"
	"container/list"
)

const (
	AC_BUFF_SIZE             = 4096
	AC_BUFF_EXHAUST_CRITICAL = 100
)

/**
 *
 * 分词器上下文状态
 *
 */
type AnalyzeContext struct {
	segmentBuff                  []rune
	charType                     []int
	bufOffset, cursor, available int
	buffLocker                   map[string]bool
	orgLexemes                   *QuickSortSet
	pathMap                      map[int]*LexemePath
	results                      *list.List
	smart                        bool
}

func NewAnalyzeContext(smart bool) (ac *AnalyzeContext) {
	ac = &AnalyzeContext{
		smart:       smart,
		segmentBuff: make([]rune, AC_BUFF_SIZE),
		charType:    make([]int, AC_BUFF_SIZE),
		buffLocker:  make(map[string]bool),
		orgLexemes:  &QuickSortSet{},
		pathMap:     make(map[int]*LexemePath),
		results:     list.New(),
		bufOffset:   0,
		cursor:      0,
		available:   0,
	}
	return
}

/**
 * 增加一个读rune的函数
 */
func readRunes(r *bufio.Reader, data []rune) (count, len int) {
	limit := cap(data)
	count = 0
	for {
		rn, width, err := r.ReadRune()
		if err != nil {
			return
		}
		data[count] = rn
		len += width
		count += 1
		if count == limit {
			return
		}
	}
}

/**
 * 根据context的上下文情况，填充segmentBuff
 * @param reader
 * @return 返回待分析的（有效的）字串长度
 */
func (ac *AnalyzeContext) fillBuffer(r *bufio.Reader) int {
	var readCount int = 0
	if ac.bufOffset == 0 {
		readCount, _ = readRunes(r, ac.segmentBuff)
	} else {
		offset := ac.available - ac.cursor - 1
		if offset > 0 {
			for i := 0; i < offset; i++ {
				ac.segmentBuff[i] = ac.segmentBuff[i+ac.cursor]
			}
			readCount = offset
		}
		// FIXED ME!
		rc, _ := readRunes(r, ac.segmentBuff[offset:])
		readCount += rc
	}
	ac.available = readCount
	ac.cursor = 0
	return readCount
}

/**
 * 初始化buff指针，处理第一个字符
 */
func (ac *AnalyzeContext) initCursor() {
	ac.cursor = 0
	ac.charType[ac.cursor] = identifyCharType(ac.segmentBuff[ac.cursor])
}

/**
 * 指针+1
 * 成功返回 true； 指针已经到了buff尾部，不能前进，返回false
 * 并处理当前字符
 */
func (ac *AnalyzeContext) moveCursor() bool {
	if ac.cursor < ac.available-1 {
		ac.cursor++
		ac.charType[ac.cursor] = identifyCharType(ac.segmentBuff[ac.cursor])
		return true
	}
	return false
}

/**
 * 设置当前segmentBuff为锁定状态
 * 加入占用segmentBuff的子分词器名称，表示占用segmentBuff
 * @param segmenterName
 */
func (ac *AnalyzeContext) lockBuffer(segmenterName string) {
	ac.buffLocker[segmenterName] = true
}

/**
 * 移除指定的子分词器名，释放对segmentBuff的占用
 * @param segmenterName
 */
func (ac *AnalyzeContext) unlockBuffer(segmenterName string) {
	if _, exists := ac.buffLocker[segmenterName]; exists {
		delete(ac.buffLocker, segmenterName)
	}
}

/**
 * 只要buffLocker中存在segmenterName
 * 则buffer被锁定
 * @return boolean 缓冲去是否被锁定
 */
func (ac *AnalyzeContext) isBufferLocked() bool {
	return len(ac.buffLocker) > 0
}

/**
 * 判断当前segmentBuff是否已经用完
 * 当前执针cursor移至segmentBuff末端this.available - 1
 * @return
 */
func (ac *AnalyzeContext) isBufferConsumed() bool {
	return ac.cursor == ac.available-1
}

/**
 * 判断segmentBuff是否需要读取新数据
 *
 * 满足一下条件时，
 * 1.available == BUFF_SIZE 表示buffer满载
 * 2.buffIndex < available - 1 && buffIndex > available - BUFF_EXHAUST_CRITICAL表示当前指针处于临界区内
 * 3.!context.isBufferLocked()表示没有segmenter在占用buffer
 * 要中断当前循环（buffer要进行移位，并再读取数据的操作）
 * @return
 */
func (ac *AnalyzeContext) needRefillBuffer() bool {
	return ac.available == AC_BUFF_SIZE &&
		ac.cursor < ac.available-1 &&
		ac.cursor > ac.available-AC_BUFF_EXHAUST_CRITICAL &&
		!ac.isBufferLocked()
}

/**
 * 累计当前的segmentBuff相对于reader起始位置的位移
 */
func (ac *AnalyzeContext) markBufferOffset() {
	ac.bufOffset += ac.cursor
}

/**
 * 向分词结果集添加词元
 * @param lexeme
 */
func (ac *AnalyzeContext) addLexeme(l *Lexeme) {
	ac.orgLexemes.addLexeme(l)
}

/**
 * 添加分词结果路径
 * 路径起始位置 ---> 路径 映射表
 * @param path
 */
func (ac *AnalyzeContext) addLexemePath(p *LexemePath) {
	if p != nil {
		ac.pathMap[p.pathBegin] = p
	}
}

/**
 * 返回原始分词结果
 * @return
 */
func (ac *AnalyzeContext) getOrgLexemes() *QuickSortSet {
	return ac.orgLexemes
}

/**
 * 对CJK字符进行单字输出
 * @param index
 */
func (ac *AnalyzeContext) outputSingleCJK(index int) {
	if CHAR_CHINESE == ac.charType[index] {
		l := NewLexeme(ac.bufOffset, index, 1, LEXEME_TYPE_CNCHAR)
		ac.results.PushBack(l)
		return
	}
	if CHAR_OTHER_CJK == ac.charType[index] {
		l := NewLexeme(ac.bufOffset, index, 1, LEXEME_TYPE_OTHER_CJK)
		ac.results.PushBack(l)
		return
	}
}

/**
 * 推送分词结果到结果集合
 * 1.从buff头部遍历到this.cursor已处理位置
 * 2.将map中存在的分词结果推入results
 * 3.将map中不存在的CJDK字符以单字方式推入results
 */
func (ac *AnalyzeContext) outputToResult() {
	var index int = 0
	for index <= ac.cursor {
		//跳过非CJK字符
		if CHAR_USELESS == ac.charType[index] {
			index++
			continue
		}

		//从pathMap找出对应index位置的LexemePath
		if p, exists := ac.pathMap[index]; exists && p != nil {
			//输出LexemePath中的lexeme到results集合
			l := p.set.pollFirst()
			for l != nil {
				ac.results.PushBack(l)
				index = l.begin + l.length
				l = p.set.pollFirst()
				if l != nil {
					for ; index < l.begin; index++ {
						ac.outputSingleCJK(index)
					}
				}
			}
		} else {
			//pathMap中找不到index对应的LexemePath
			//单字输出
			ac.outputSingleCJK(index)
			index++
		}
	}
	ac.pathMap = make(map[int]*LexemePath)
}

/**
 * 组合词元
 */
func (ac *AnalyzeContext) compound(l *Lexeme) {
	if !ac.smart {
		return
	}

	//数量词合并处理
	if ac.results.Len() != 0 {
		if LEXEME_TYPE_ARABIC == l.lexemeType {
			n := ac.results.Front().Value.(*Lexeme)
			appendOK := false
			if LEXEME_TYPE_CNUM == n.lexemeType {
				//合并英文数词+中文数词
				appendOK = l.append(n, LEXEME_TYPE_CNUM)
			} else if LEXEME_TYPE_COUNT == n.lexemeType {
				//合并英文数词+中文量词
				appendOK = l.append(n, LEXEME_TYPE_CQUAN)
			}
			if appendOK {
				ac.results.Remove(ac.results.Front())
			}
		}
	}

	//可能存在第二轮合并
	if LEXEME_TYPE_CNUM == l.lexemeType && ac.results.Len() != 0 {
		n := ac.results.Front().Value.(*Lexeme)
		appendOK := false
		if LEXEME_TYPE_COUNT == n.lexemeType {
			//合并中文数词+中文量词
			appendOK = l.append(n, LEXEME_TYPE_CQUAN)
		}
		if appendOK {
			ac.results.Remove(ac.results.Front())
		}
	}
}

/**
 * 返回lexeme
 *
 * 同时处理合并
 * @return
 */
func (ac *AnalyzeContext) getNextLexeme() (l *Lexeme) {
	//从结果集取出，并移除第一个Lexme
	el := ac.results.Front()
	if el == nil {
		l = nil
		return
	}
	result := el.Value.(*Lexeme)
	ac.results.Remove(el)

	for result != nil {
		ac.compound(result)
		if isStopWord(ac.segmentBuff, result.begin, result.length) {
			//是停止词继续取列表的下一个
			el := ac.results.Front()
			if el == nil {
				l = nil
				return
			}
			result = el.Value.(*Lexeme)
			ac.results.Remove(el)
		} else {
			//不是停止词, 生成lexeme的词元文本,输出
			result.lexemeText = string(ac.segmentBuff[result.begin : result.begin+result.length])
			break
		}
	}
	l = result
	return
}

/**
 * 重置分词上下文状态
 */
func (ac *AnalyzeContext) reset() {
	ac.buffLocker = make(map[string]bool)
	ac.orgLexemes = &QuickSortSet{}
	ac.available = 0
	ac.bufOffset = 0
	ac.charType = make([]int, AC_BUFF_SIZE)
	ac.cursor = 0
	ac.results = list.New()
	ac.segmentBuff = make([]rune, AC_BUFF_SIZE)
	ac.pathMap = make(map[int]*LexemePath)
}
