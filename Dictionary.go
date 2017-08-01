package ikgo

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var (
	MainDict, SurnameDict, QuantifierDict, SuffixDict, PrepDict, StopWords *DictSegment
	conf_dir                                                               string
	ext_files, ext_stopfiles                                               []string
	conf_smart                                                             bool
)

const (
	PATH_DIC_MAIN       = "main.dic"
	PATH_DIC_SURNAME    = "surname.dic"
	PATH_DIC_QUANTIFIER = "quantifier.dic"
	PATH_DIC_SUFFIX     = "suffix.dic"
	PATH_DIC_PREP       = "preposition.dic"
	PATH_DIC_STOP       = "stopword.dic"

	FILE_NAME = "IKAnalyzer.cfg.xml"
	EXT_DICT  = "ext_dict"
	EXT_STOP  = "ext_stopwords"

	//REMOTE_EXT_DICT = "remote_ext_dict"
	//REMOTE_EXT_STOP = "remote_ext_stopwords"
)

type DictFile struct {
	Key  string `xml:"key,attr"`
	Path string `xml:",chardata"`
}

type DictFiles struct {
	DictFiles []DictFile `xml:"entry"`
}

func InitDict(dir string, bs bool) {
	conf_dir = dir
	conf_smart = bs
	// 读取cfg.xml
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", conf_dir, FILE_NAME))
	if err == nil {
		dfs := &DictFiles{}
		err = xml.Unmarshal(content, dfs)
		if err == nil {
			for _, df := range dfs.DictFiles {
				if df.Key == EXT_DICT {
					ext_files = strings.Split(df.Path, ";")
				} else if df.Key == EXT_STOP {
					ext_stopfiles = strings.Split(df.Path, ";")
				}

			}
		}

	}

	loadMainDict()
	loadSurnameDict()
	loadQuantifierDict()
	loadSuffixDict()
	loadPrepDict()
	loadStopWordDict()
}

/**
 * 加载主词典及扩展词典
 */
func loadMainDict() {
	MainDict = NewDictSegment(0)

	// 读取主词典文件
	fi, err := os.Open(fmt.Sprintf("%s/%s", conf_dir, PATH_DIC_MAIN))
	if err != nil {
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		word, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if word != nil {
			trimWord := bytes.Trim(word, "\r\n\t ")
			if !bytes.Equal(trimWord, []byte("")) {
				MainDict.fillSegment([]rune(string(trimWord)))
			}

		}
	}

	// 加载扩展词典
	loadExtDict()
}

/**
 * 加载用户配置的扩展词典到主词库表
 */
func loadExtDict() {
	// 加载扩展词典配置
	for _, fname := range ext_files {
		fi, err := os.Open(fmt.Sprintf("%s/%s", conf_dir, fname))
		if err != nil {
			continue
		}
		br := bufio.NewReader(fi)
		for {
			word, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			if word != nil {
				trimWord := bytes.Trim(word, "\r\n\t ")
				if !bytes.Equal(trimWord, []byte("")) {
					MainDict.fillSegment([]rune(string(trimWord)))
				}

			}
		}
		fi.Close()
	}
}

/**
 * 加载用户扩展的停止词词典
 */
func loadStopWordDict() {
	StopWords = NewDictSegment(0)

	// 读取主词典文件
	fi, err := os.Open(fmt.Sprintf("%s/%s", conf_dir, PATH_DIC_STOP))
	if err != nil {
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		word, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if word != nil {
			trimWord := bytes.Trim(word, "\r\n\t ")
			if !bytes.Equal(trimWord, []byte("")) {
				StopWords.fillSegment([]rune(string(trimWord)))
			}

		}
	}

	// 加载扩展词典
	// 加载扩展词典配置
	for _, fname := range ext_stopfiles {
		fi, err := os.Open(fmt.Sprintf("%s/%s", conf_dir, fname))
		if err != nil {
			continue
		}
		br := bufio.NewReader(fi)
		for {
			word, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			if word != nil {
				trimWord := bytes.Trim(word, "\r\n\t ")
				if !bytes.Equal(trimWord, []byte("")) {
					MainDict.fillSegment([]rune(string(trimWord)))
				}

			}
		}
		fi.Close()
	}
}

/**
 * 加载量词词典
 */
func loadQuantifierDict() {
	// 建立一个量词典实例
	QuantifierDict = NewDictSegment(0)

	// 读取主词典文件
	fi, err := os.Open(fmt.Sprintf("%s/%s", conf_dir, PATH_DIC_QUANTIFIER))
	if err != nil {
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		word, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if word != nil {
			trimWord := bytes.Trim(word, "\r\n\t ")
			if !bytes.Equal(trimWord, []byte("")) {
				QuantifierDict.fillSegment([]rune(string(trimWord)))
			}

		}
	}
}

/**
 * 加载量词词典
 */
func loadSurnameDict() {
	// 建立一个量词典实例
	SurnameDict = NewDictSegment(0)

	// 读取主词典文件
	fi, err := os.Open(fmt.Sprintf("%s/%s", conf_dir, PATH_DIC_SURNAME))
	if err != nil {
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		word, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if word != nil {
			trimWord := bytes.Trim(word, "\r\n\t ")
			if !bytes.Equal(trimWord, []byte("")) {
				SurnameDict.fillSegment([]rune(string(trimWord)))
			}

		}
	}
}

func loadSuffixDict() {
	// 建立一个量词典实例
	SuffixDict = NewDictSegment(0)

	// 读取主词典文件
	fi, err := os.Open(fmt.Sprintf("%s/%s", conf_dir, PATH_DIC_SUFFIX))
	if err != nil {
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		word, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if word != nil {
			trimWord := bytes.Trim(word, "\r\n\t ")
			if !bytes.Equal(trimWord, []byte("")) {
				SuffixDict.fillSegment([]rune(string(trimWord)))
			}

		}
	}
}

func loadPrepDict() {
	// 建立一个量词典实例
	PrepDict = NewDictSegment(0)

	// 读取主词典文件
	fi, err := os.Open(fmt.Sprintf("%s/%s", conf_dir, PATH_DIC_PREP))
	if err != nil {
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		word, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if word != nil {
			trimWord := bytes.Trim(word, "\r\n\t ")
			if !bytes.Equal(trimWord, []byte("")) {
				PrepDict.fillSegment([]rune(string(trimWord)))
			}

		}
	}
}

/**
 * 判断是否是停止词
 *
 * @return boolean
 */
func isStopWord(charArray []rune, begin, length int) bool {
	return StopWords.matchSeg(charArray, begin, length).isMatch()
}

/**
 * 从已匹配的Hit中直接取出DictSegment，继续向下匹配
 *
 * @return Hit
 */
func matchWithHit(charArray []rune, currentIndex int, matchedHit *Hit) *Hit {
	ds := matchedHit.matchedDictSegment
	return ds.matchSegSearch(charArray, currentIndex, 1, matchedHit)
}
