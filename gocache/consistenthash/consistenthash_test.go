package consistenthash

import (
	"strconv"
	"testing"
)

// 测试 hash
func TestHashing(t *testing.T) {
	// 使用自定义的 hash 函数（返回对应的数字）方便测试
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26 哈希环
	hash.Add("6", "4", "2")

	// 用例及对应的真实节点
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8, 18, 28
	hash.Add("8")

	// 添加一个真实节点 8, 对应虚拟节点的哈希值是 08/18/28, 此时，
	// 用例 27 对应的虚拟节点从 02 变更为 28, 即真实节点 8
	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}
