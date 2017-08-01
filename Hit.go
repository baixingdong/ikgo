package ikgo

const (
	HIT_UNMATCH = 0x00000000
	HIT_MATCH   = 0x00000001
	HIT_PREFIX  = 0x00000010
)

type Hit struct {
	hitState           int          //该HIT当前状态，默认未匹配
	matchedDictSegment *DictSegment //记录词典匹配过程中，当前匹配到的词典分支节点
	beg, end           int          //词段起止位置
}

/**
 * 判断是否完全匹配
 */
func (h *Hit) isMatch() bool {
	return (h.hitState & HIT_MATCH) > 0
}

func (h *Hit) setMatch() {
	h.hitState |= HIT_MATCH
}

/**
 * 判断是否是词的前缀
 */

func (h *Hit) isPrefix() bool {
	return (h.hitState & HIT_PREFIX) > 0
}

func (h *Hit) setPrefix() {
	h.hitState |= HIT_PREFIX
}

/**
 * 判断是否是不匹配
 */
func (h *Hit) isUnmatch() bool {
	return h.hitState == HIT_UNMATCH
}

func (h *Hit) setUnmatch() {
	h.hitState = HIT_UNMATCH
}
