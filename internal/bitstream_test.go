package internal

import (
	"bufio"
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestScanNalu(t *testing.T) {
	tokens := [][]byte{
		{},
		{1},
		{1, 2},
		{1, 0, 3},
		{1, 0, 1, 4},
		{1, 0, 1, 4, 5},
	}
	indexs := make([]int, len(tokens))
	for i := 0; i < cap(indexs); i++ {
		indexs[i] = i
	}
	Shuffle(indexs)
	t.Logf("indexs = %v", indexs)

	//bitStream := make([]byte, 0)
	bitStream := annexBSpliter1
	for i := 0; i < len(indexs); i++ {
		bitStream = append(bitStream, tokens[indexs[i]]...)
		if i%3 == 0 {
			bitStream = append(bitStream, annexBSpliter1...)
		} else {
			bitStream = append(bitStream, annexBSpliter2...)
		}
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(bitStream))
	scanner.Split(ScanNalu)

	for i := 0; i < len(indexs);{
		if len(tokens[indexs[i]]) < 1 {
			t.Logf("skip token idx = %v, value = %v", indexs[i], tokens[indexs[i]])
			i++
			continue
		}
		if !scanner.Scan() {
			break
		}
		if len(scanner.Bytes()) < 1 {
			continue
		}
		t.Logf("got token %v of index %v", scanner.Bytes(), indexs[i])
		if !bytes.Equal(tokens[indexs[i]], scanner.Bytes()) {
			t.Errorf("token %v not equal orig:%v, got:%v", indexs[i], tokens[indexs[i]], scanner.Bytes())
		}
		i++
	}
	//i := 0
	//for scanner.Scan() {
	//	t.Logf("got token %v of index %v", scanner.Bytes(), indexs[i])
	//	if len(scanner.Bytes()) < 1 {
	//		i++
	//		continue
	//	}
	//	if !bytes.Equal(tokens[indexs[i]], scanner.Bytes()) {
	//		t.Errorf("token %v not equal orig:%v, got:%v", indexs[i], tokens[indexs[i]], scanner.Bytes())
	//	}
	//	i++
	//}
	//if i < len(tokens) {
	//	t.Errorf("input token count = %v, output token count = %v", len(tokens), i)
	//}

	if scanner.Err() != nil {
		t.Errorf("scanner error = %v", scanner.Err())
	}
}

func Shuffle(vals []int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}
