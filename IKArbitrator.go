package ikgo

import (
	"container/list"
)

type IKArbitrator struct {
}

/**
 * 向前遍历，添加词元，构造一个无歧义词元组合
 * @param LexemePath path
 * @return
 */
func (a *IKArbitrator) forwardPath(lexemeCell *Cell, option *LexemePath) *list.List {
	//发生冲突的Lexeme栈
	conflictStack := list.New()
	c := lexemeCell
	//迭代遍历Lexeme链表
	for c != nil && c.lexeme != nil {
		if !option.addNotCrossLexeme(c.lexeme) {
			//词元交叉，添加失败则加入lexemeStack栈
			conflictStack.PushBack(c)
		}
		c = c.next
	}
	return conflictStack
}

/**
 * 回滚词元链，直到它能够接受指定的词元
 * @param lexeme
 * @param l
 */
func (a *IKArbitrator) backPath(l *Lexeme, option *LexemePath) {
	for option.checkCross(l) {
		option.removeTail()
	}
}

/**
 * 歧义识别
 * @param lexemeCell 歧义路径链表头
 * @param fullTextLength 歧义路径文本长度
 * @return
 */
func (a *IKArbitrator) judge(lexemeCell *Cell, fullTextLength int) *LexemePath {
	//候选路径集合
	pathOptions := []*LexemePath{}
	//候选结果路径
	option := NewLexemePath()
	//对crossPath进行一次遍历,同时返回本次遍历中有冲突的Lexeme栈
	lexemeStack := a.forwardPath(lexemeCell, option)

	//当前词元链并非最理想的，加入候选路径集合
	pathOptions = append(pathOptions, option.deepCopy())

	//存在歧义词，处理
	var c *Cell = nil
	for lexemeStack.Len() != 0 {
		el := lexemeStack.Back()
		c = el.Value.(*Cell)
		lexemeStack.Remove(el)
		//回滚词元链
		a.backPath(c.lexeme, option)
		//从歧义词位置开始，递归，生成可选方案
		a.forwardPath(c, option)

		pathOptions = append(pathOptions, option.deepCopy())
	}

	//返回集合中的最优方案
	best := pathOptions[0]
	for _, po := range pathOptions {
		if best.compare(po) > 0 {
			best = po
		}
	}
	return best
}

/**
 * 分词歧义处理
 * @param orgLexemes
 * @param useSmart
 */
func (a *IKArbitrator) process(context *AnalyzeContext, useSmart bool) {
	orgLexemes := context.getOrgLexemes()
	orgLexeme := orgLexemes.pollFirst()

	crossPath := NewLexemePath()
	for orgLexeme != nil {
		if !crossPath.addCrossLexeme(orgLexeme) {
			//找到与crossPath不相交的下一个crossPath
			if crossPath.set.size == 1 || !useSmart {
				//crossPath没有歧义 或者 不做歧义处理
				//直接输出当前crossPath
				context.addLexemePath(crossPath)
			} else {
				//对当前的crossPath进行歧义处理
				headCell := crossPath.set.head
				judgeResult := a.judge(headCell, crossPath.getPathLength())
				//输出歧义处理结果judgeResult
				context.addLexemePath(judgeResult)
			}

			//把orgLexeme加入新的crossPath中
			crossPath = NewLexemePath()
			crossPath.addCrossLexeme(orgLexeme)
		}

		orgLexeme = orgLexemes.pollFirst()
	}

	//处理最后的path
	if crossPath.set.size == 1 || !useSmart {
		//crossPath没有歧义 或者 不做歧义处理
		//直接输出当前crossPath
		context.addLexemePath(crossPath)
	} else {
		//对当前的crossPath进行歧义处理
		headCell := crossPath.set.head
		judgeResult := a.judge(headCell, crossPath.getPathLength())
		//输出歧义处理结果judgeResult
		context.addLexemePath(judgeResult)
	}
}
