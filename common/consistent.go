package common

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type units []uint32

func (x units) Len() int {
	return len(x)
}

func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

// 交换值
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// 创建结构体保存一致性hash信息
type Consistent struct {
	// hash环，key为哈希值，值存放节点信息
	circle         map[uint32]string
	// 已经排序的节点hash切片
	sortedHashes   units
	// 虚拟节点个数，增加hash的平衡性
	VirtualNodeNum int
	sync.RWMutex
}

func NewConsistent() *Consistent {
	consistent := &Consistent{
		circle:         make(map[uint32]string),
		VirtualNodeNum: 20,
	}
	return consistent
}
// 生成虚拟节点的key
func (c *Consistent) generateKey(element string, index int) string{
	return element + strconv.Itoa(index)
}
// 根据key生成hash值
func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		// 声明一个数组长度为64
		var srcatch [64]byte
		copy(srcatch[:], key)
		// 使用IEEE多项式返回数据的CRC-32校验和
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	// 判断容量是否过大，过大则重置
	standardCap := (c.VirtualNodeNum*4) *len(c.circle)
	if cap(c.sortedHashes) > standardCap {
		hashes = make([]uint32, standardCap)
	}
	// 添加hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	// 对所有节点hash值进行排序，方便二分查找
	sort.Sort(hashes)
	c.sortedHashes = hashes
}
// 添加节点
func (c *Consistent) add(element string) {
	for i:=0;i<c.VirtualNodeNum;i++ {
		key := c.generateKey(element, i)
		hashedValueByKey := c.hashKey(key) // 根据element和i生成的key来生成hash值
		c.circle[hashedValueByKey] = element// 将hash值放在圆环上
		_ =append(c.sortedHashes, hashedValueByKey)
		c.updateSortedHashes()
	}
}

func (c *Consistent) Add(element string) {
	c.Lock() // 写操作加写锁，读和写都阻塞
	defer c.Unlock()
	c.add(element)
}


func (c *Consistent) remove(element string) {
	for i := 0;i < c.VirtualNodeNum; i ++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

func(c *Consistent) search(key uint32) int {
	// 查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	// 使用二分法来搜索指定切片满足条件的最小值
	i := sort.Search(len(c.sortedHashes), f)
	// 如果超出范围则设置i=0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

func (c *Consistent) Get(name string) (string, error) {
	c.RLock() // 读锁
	defer c.RUnlock()
	if len(c.circle) == 0{
		return "", nil
	}
	key := c.hashKey(name)
	i := c.search(key)
	//fmt.Println(c.circle)
	return c.circle[c.sortedHashes[i]], nil
}