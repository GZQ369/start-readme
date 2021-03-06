package model

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

const (
	sdsHdr     = "sdsHdr"
	sdsInt     = "sdsInt"
	zipList    = "ziplist"
	linkedList = "linkedlist"
	intSet = "IntSet"
	dictSet = "dictSet"
)

type redisObject struct {
	Type      interface{}
	Enconding interface{} //编码类型
	lru       int64       //记录对象最后一次被应用程序访问的时间，和当前时间的差值可得出空转的时间
	Refound   int         //记录对象被应用的次数，当初始化时次数为1，被引用后加一，当次数为0时，内存被回收
	ptr       unsafe.Pointer
}
type String struct{}
type List struct{}
type Set struct{}
type Hash struct{}
type Zset struct{}

type RedisDb struct {
	dict    map[string]redisObject
	expires map[string]int64 //过期时间
}

func RedisNew() RedisDb {
	return RedisDb{dict: make(map[string]redisObject)}
}
func stringObjectNew(v string) (redisObject, error) {
	va, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return redisObject{Type: String{}, Enconding: sdsHdr, lru: time.Now().Unix(), Refound: 1, ptr: unsafe.Pointer(sdsHdrNew(v))}, nil
	}
	return redisObject{Type: String{}, Enconding: sdsInt, lru: time.Now().Unix(), Refound: 1, ptr: unsafe.Pointer(SdsIntNew(va))}, nil

}
func dictObjectNew(dc *dictht) redisObject {
	return redisObject{Type: Hash{}, Enconding: dictht{}, lru: time.Now().Unix(), Refound: 0, ptr: unsafe.Pointer(dc)}
}
func listObjectNew() (redisObject, error) {
	return redisObject{Type: List{}, Enconding: ChainList{}, lru: time.Now().Unix(), Refound: 0, ptr: unsafe.Pointer(ChainListCreate())}, nil
}
func zsetObjectNew() redisObject {
	return redisObject{Type: Zset{}, Enconding: ZSkipList{}, lru: time.Now().Unix(), Refound: 0, ptr: unsafe.Pointer(new(ZSkipList))}
}
func setObjectNew(v interface{}) redisObject {

	var res unsafe.Pointer
	var strTmp string
	switch x := v.(type) {
	case []string:
		fmt.Printf("x is a string，value is %v\n", v)
		res = unsafe.Pointer(IntSetNew(x))
		strTmp = intSet
	case *dictht:
		fmt.Printf("x is a int is %v\n", v)
		res = unsafe.Pointer(new(dictht))
		strTmp = dictSet
	}
	return redisObject{Type: String{}, Enconding: strTmp, lru: time.Now().Unix(), Refound: 1, ptr: res}

}
//返回现在所有对象中，各个对象的数量
func (r RedisDb) GetObjectNum() map[interface{}]int64 {
	res := map[interface{}]int64{}
	for _, ob := range r.dict {
		res[ob.Type] ++
	}
	return res
}
