package main

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strings"

	"github.com/go-ego/gse"
)

// 初始化中文分词器
var segmenter gse.Segmenter

func init() {
	_ = segmenter.LoadDict()
}

// 判断文本是否包含中文、日文、韩文字符
func containsCJK(text string) bool {
	cjkRegex := regexp.MustCompile(`[\p{Han}\p{Hiragana}\p{Katakana}\p{Hangul}]`)
	return cjkRegex.MatchString(text)
}

// tokenize 分词(简单分词)
func tokenize(text string) []string {
	if containsCJK(text) {
		return segmenter.Cut(text, true) // 中文、日文、韩文使用分词
	}

	return strings.Fields(text) // 其他语言使用空格分词
}

// hashToken 计算单个token的哈希值
func hashToken(token string) uint64 {
	h64 := fnv.New64a()
	_, _ = h64.Write([]byte(token))
	return h64.Sum64()
}

// simHash 计算SimHash指纹
func simHash(text string) uint64 {
	tokens := tokenize(text)
	if len(tokens) == 0 {
		return 0 // 避免空文本计算
	}

	var vector [64]int
	for _, token := range tokens {
		hash := hashToken(token)
		for i := 0; i < 64; i++ {
			if (hash>>i)&1 == 1 {
				vector[i]++
			} else {
				vector[i]--
			}
		}
	}

	var fingerprint uint64
	for i := 0; i < 64; i++ {
		if vector[i] > 0 {
			fingerprint |= 1 << i
		}
	}
	return fingerprint
}

// hammingDistance 计算汉明距离
func hammingDistance(hash1, hash2 uint64) int {
	xor := hash1 ^ hash2
	distance := 0
	for xor > 0 {
		distance += int(xor & 1)
		xor >>= 1
	}
	return distance
}

// similarity 计算相似度
func similarity(hash1, hash2 uint64) float64 {
	if hash1 == 0 || hash2 == 0 {
		return 0.0 // 避免短文本导致相似度过高
	}
	distance := hammingDistance(hash1, hash2)
	return 1 - float64(distance)/64.0
}

func main() {
	text1 := `This was a five star "Good morning sir." She curtsied and smiled sheepishly.`
	text2 := `She was a beautiful girl with an unusual personality. Coming late to work and lying her way out was her thing.`

	/*
		常见的相似度范围指标：
			0.9 以上：非常相似，几乎是相同文本，可能只有少量的不同词汇。
			0.7 ~ 0.9：相似，可能是相同主题的文本，修改了一些词语或句子结构。
			0.5 ~ 0.7：中等相似度，可能是不同的文本，但包含了一些相似的信息。
			低于 0.5：不相似，文本内容差异较大。
	*/
	hash1 := simHash(text1)
	hash2 := simHash(text2)
	fmt.Printf("文本1的 SimHash: %064b\n", hash1)
	fmt.Printf("文本2的 SimHash: %064b\n", hash2)
	fmt.Printf("汉明距离: %d\n", hammingDistance(hash1, hash2)) // 通常汉明距离较小，文本非常相似。相似度=1−(汉明距离/64)
	fmt.Printf("相似度: %.4f\n", similarity(hash1, hash2))     // 0 - 1 之间。越接近1，表示两个文本越相似，反之。
}
