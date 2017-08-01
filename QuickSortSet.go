package ikgo

/**
 * IK分词器专用的Lexem快速排序集合
 */

type Cell struct {
	lexeme     *Lexeme
	prev, next *Cell
}

type QuickSortSet struct {
	head, tail *Cell
	size       int
}

func (q *QuickSortSet) addLexeme(l *Lexeme) bool {
	ne := &Cell{lexeme: l, prev: nil, next: nil}
	if q.size == 0 {
		q.tail = ne
		q.head = q.tail
		q.size++
		return true
	}

	if q.tail.lexeme.compare(l) == 0 {
		return false
	}
	if q.tail.lexeme.compare(l) < 0 {
		q.tail.next = ne
		ne.prev = q.tail
		q.tail = ne
		q.size++
		return true
	}

	if q.head.lexeme.compare(l) > 0 {
		q.head.prev = ne
		ne.next = q.head
		q.head = ne
		q.size++
		return true
	}

	// 从尾部逆上
	index := q.tail
	for index != nil && index.lexeme.compare(l) > 0 {
		index = index.prev
	}
	if index.lexeme.compare(l) == 0 {
		return false
	}
	if index.lexeme.compare(l) < 0 {
		ne.prev = index
		ne.next = index.next
		index.next.prev = ne
		index.next = ne
		q.size++
		return true
	}
	return false
}

/**
 * 返回链表头部元素
 * @return
 */
func (q *QuickSortSet) peekFirst() *Lexeme {
	if q.head != nil {
		return q.head.lexeme
	}
	return nil
}

/**
 * 取出链表集合的第一个元素
 * @return Lexeme
 */
func (q *QuickSortSet) pollFirst() (l *Lexeme) {
	if q.size > 0 {
		head := q.head
		l = head.lexeme
		q.head = head.next
		if q.head == nil {
			q.tail = nil
		}
		q.size--
		return
	}
	l = nil
	return
}

/**
 * 返回链表尾部元素
 * @return
 */
func (q *QuickSortSet) peekLast() *Lexeme {
	if q.tail != nil {
		return q.tail.lexeme
	}
	return nil
}

/**
 * 取出链表集合的最后一个元素
 * @return Lexeme
 */
func (q *QuickSortSet) pollLast() (l *Lexeme) {
	if q.size > 0 {
		tail := q.tail
		l = tail.lexeme
		q.tail = tail.prev
		if q.tail == nil {
			q.head = nil
		}
		q.size--
		return
	}
	l = nil
	return
}
