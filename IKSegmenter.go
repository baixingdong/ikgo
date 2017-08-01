package ikgo

import (
	"bufio"
	"strings"
)

// 重写ik分词

type IKSegmenter struct {
	reader     *bufio.Reader
	context    *AnalyzeContext
	segmenters []ISegmenter
	arbitrator IKArbitrator
	useSmart   bool
}

func init() {
	initCNQS()
	initLS()
}

func NewIKSegmenter(input string, useSmart bool) *IKSegmenter {
	ret := &IKSegmenter{
		reader:     bufio.NewReader(strings.NewReader(input)),
		context:    NewAnalyzeContext(useSmart),
		arbitrator: IKArbitrator{},
		useSmart:   useSmart,
	}
	ret.loadSegmenters()
	return ret
}

/**
 * 初始化词典，加载子分词器实现
 * @return List<ISegmenter>
 */
func (s *IKSegmenter) loadSegmenters() {
	s.segmenters = []ISegmenter{
		NewLetterSegmenter(),
		NewCN_QuantifierSegmenter(),
		NewCJKSegmenter(),
	}
}

/**
 * 分词，获取下一个词元
 * @return Lexeme 词元对象
 * @throws java.io.IOException
 */
func (s *IKSegmenter) Next() *Lexeme {
	var l *Lexeme = s.context.getNextLexeme()
	for l == nil {
		/*
		 * 从reader中读取数据，填充buffer
		 * 如果reader是分次读入buffer的，那么buffer要  进行移位处理
		 * 移位处理上次读入的但未处理的数据
		 */
		available := s.context.fillBuffer(s.reader)
		if available <= 0 {
			//reader已经读完
			s.context.reset()
			return nil
		}

		//初始化指针
		s.context.initCursor()
		for {
			//遍历子分词器
			for _, segmenter := range s.segmenters {
				segmenter.analyze(s.context)
			}

			//字符缓冲区接近读完，需要读入新的字符
			if s.context.needRefillBuffer() {
				break
			}

			//向前移动指针
			if !s.context.moveCursor() {
				break
			}
		}
		//重置子分词器，为下轮循环进行初始化
		for _, segmenter := range s.segmenters {
			segmenter.reset()
		}

		//对分词进行歧义处理
		s.arbitrator.process(s.context, s.useSmart)
		//将分词结果输出到结果集，并处理未切分的单个CJK字符
		s.context.outputToResult()
		//记录本次分词的缓冲区位移
		s.context.markBufferOffset()

		l = s.context.getNextLexeme()
	}
	return l
}

/**
 * 重置分词器到初始状态
 * @param input
 */
func (s *IKSegmenter) Reset(input string) {
	s.context.reset()
	for _, segmenter := range s.segmenters {
		segmenter.reset()
	}
	s.reader = bufio.NewReader(strings.NewReader(input))
}
