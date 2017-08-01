package ikgo

type ISegmenter interface {
	/**
	 * 从分析器读取下一个可能分解的词元对象
	 * @param context 分词算法上下文
	 */
	analyze(context *AnalyzeContext)

	/**
	 * 重置子分析器状态
	 */
	reset()
}
