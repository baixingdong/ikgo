package ikgo

import (
	"sort"
)

var (
	charMap            map[rune]rune = make(map[rune]rune)
	ARRAY_LENGTH_LIMIT int           = 3
)

type DictSegment struct {
	childrenMap   map[rune]*DictSegment //Map存储结构
	childrenArray []*DictSegment        //数组方式存储结构
	nodeChar      rune                  //当前节点上存储的字符
	storeSize     int                   //当前节点存储的Segment数目 ==> storeSize <=ARRAY_LENGTH_LIMIT ，使用数组存储， storeSize >ARRAY_LENGTH_LIMIT ,则使用Map存储
	nodeState     int                   //当前DictSegment状态 ,默认 0 , 1表示从根节点到当前节点的路径表示一个词
}

func NewDictSegment(nodeChar rune) *DictSegment {
	return &DictSegment{nodeChar: nodeChar}
}

/*
 * 判断是否有下一个节点
 */
func (ds *DictSegment) hasNextNode() bool {
	return ds.storeSize > 0
}

/**
 * 匹配词段
 * @param charArray
 * @param begin
 * @param length
 * @param searchHit
 * @return Hit
 */
func (ds *DictSegment) matchSegSearch(charArray []rune, begin, length int, searchHit *Hit) *Hit {
	if searchHit == nil {
		//如果hit为空，新建
		searchHit = &Hit{}
		searchHit.beg = begin
	} else {
		//否则要将HIT状态重置
		searchHit.setUnmatch()
	}
	//设置hit的当前处理位置
	searchHit.end = begin
	keyChar := charArray[begin]

	//引用实例变量为本地变量，避免查询时遇到更新的同步问题
	segmentArray := ds.childrenArray
	segmentMap := ds.childrenMap
	var nds *DictSegment = nil

	//STEP1 在节点中查找keyChar对应的DictSegment
	if len(segmentArray) > 0 {
		keySegment := NewDictSegment(keyChar)
		position := sort.Search(ds.storeSize, func(i int) bool {
			return keySegment.nodeChar == segmentArray[i].nodeChar
		})
		if position != ds.storeSize {
			nds = segmentArray[position]
		}
	} else if len(segmentMap) > 0 {
		nds, _ = segmentMap[keyChar]
	}

	//STEP2 找到DictSegment，判断词的匹配状态，是否继续递归，还是返回结果
	if nds != nil {
		if length > 1 {
			//词未匹配完，继续往下搜索
			return nds.matchSegSearch(charArray, begin+1, length-1, searchHit)
		}
		if length == 1 {
			//搜索最后一个char
			if nds.nodeState == 1 {
				//添加HIT状态为完全匹配
				searchHit.setMatch()
			}
			if nds.hasNextNode() {
				//添加HIT状态为前缀匹配
				searchHit.setPrefix()
				//记录当前位置的DictSegment
				searchHit.matchedDictSegment = nds
			}
			return searchHit
		}
	}
	//STEP3 没有找到DictSegment， 将HIT设置为不匹配
	return searchHit
}

/**
 * 匹配词段
 * @param charArray
 * @param begin
 * @param length
 * @return Hit
 */
func (ds *DictSegment) matchSeg(charArray []rune, begin, length int) *Hit {
	return ds.matchSegSearch(charArray, begin, length, nil)
}

/**
 * 匹配词段
 * @param charArray
 * @return Hit
 */
func (ds *DictSegment) match(charArray []rune) *Hit {
	return ds.matchSeg(charArray, 0, len(charArray))
}

func (ds *DictSegment) getChildrenArray() []*DictSegment {
	if ds.childrenArray == nil {
		ds.childrenArray = make([]*DictSegment, ARRAY_LENGTH_LIMIT)
	}
	return ds.childrenArray
}

func (ds *DictSegment) getChildrenMap() map[rune]*DictSegment {
	if ds.childrenMap == nil {
		ds.childrenMap = make(map[rune]*DictSegment, ARRAY_LENGTH_LIMIT*2)
	}
	return ds.childrenMap
}

/**
 * 查找本节点下对应的keyChar的segment	 *
 * @param keyChar
 * @param create  =1如果没有找到，则创建新的segment ; =0如果没有找到，不创建，返回null
 * @return
 */
func (ds *DictSegment) lookforSegment(keyChar rune, create int) *DictSegment {
	var nds *DictSegment = nil
	if ds.storeSize <= ARRAY_LENGTH_LIMIT {
		//获取数组容器，如果数组未创建则创建数组
		segmentArray := ds.getChildrenArray()
		//搜寻数组
		keySegment := NewDictSegment(keyChar)
		position := sort.Search(ds.storeSize, func(i int) bool {
			return keySegment == segmentArray[i]
		})
		if position != len(segmentArray) {
			nds = segmentArray[position]
		}

		//遍历数组后没有找到对应的segment
		if nds == nil && create == 1 {
			nds = keySegment
			if ds.storeSize < ARRAY_LENGTH_LIMIT {
				//数组容量未满，使用数组存储，插入
				idx := ds.storeSize
				for i := ds.storeSize; i > 0; i-- {
					if segmentArray[i-1].nodeChar < nds.nodeChar {
						idx = i
						break
					} else {
						segmentArray[i] = segmentArray[i-1]
					}
				}
				segmentArray[idx] = nds
				//segment数目+1
				ds.storeSize++
			} else {
				//数组容量已满，切换Map存储
				//获取Map容器，如果Map未创建,则创建Map
				segmentMap := ds.getChildrenMap()
				//将数组中的segment迁移到Map中
				for i := 0; i < ds.storeSize; i++ {
					segmentMap[segmentArray[i].nodeChar] = segmentArray[i]
				}
				//存储新的segment
				segmentMap[keyChar] = nds
				//segment数目+1 ，  必须在释放数组前执行storeSize++ ， 确保极端情况下，不会取到空的数组
				ds.storeSize++
				//释放当前的数组引用
				ds.childrenArray = nil
			}
		}
	} else {
		//获取Map容器，如果Map未创建,则创建Map
		segmentMap := ds.getChildrenMap()
		//搜索Map
		var exists bool
		nds, exists = segmentMap[keyChar]
		if !exists && create == 1 {
			//构造新的segment
			nds = NewDictSegment(keyChar)
			segmentMap[keyChar] = nds
			//当前节点存储segment数目+1
			ds.storeSize++
		}
	}
	return nds
}

/**
 * 加载填充词典片段
 * @param charArray
 * @param begin
 * @param length
 * @param enabled
 */
func (ds *DictSegment) fillSegmentSeg(charArray []rune, begin, length, enabled int) {
	//获取字典表中的汉字对象
	beginChar := charArray[begin]
	keyChar, exists := charMap[beginChar]
	//字典中没有该字，则将其添加入字典
	if !exists {
		charMap[beginChar] = beginChar
		keyChar = beginChar
	}

	//搜索当前节点的存储，查询对应keyChar的keyChar，如果没有则创建
	nds := ds.lookforSegment(keyChar, enabled)
	if nds != nil {
		//处理keyChar对应的segment
		if length > 1 {
			//词元还没有完全加入词典树
			nds.fillSegmentSeg(charArray, begin+1, length-1, enabled)
		} else if length == 1 {
			//已经是词元的最后一个char,设置当前节点状态为enabled，
			//enabled=1表明一个完整的词，enabled=0表示从词典中屏蔽当前词
			nds.nodeState = enabled
		}
	}

}

/**
 * 加载填充词典片段
 * @param charArray
 */
func (ds *DictSegment) fillSegment(charArray []rune) {
	ds.fillSegmentSeg(charArray, 0, len(charArray), 1)
}
