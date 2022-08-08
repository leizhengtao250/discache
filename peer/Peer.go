package peer

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool) //根据相应的key获取对应的节点
}

type PeerGetter interface {
	Get(group string, key string) ([]byte, error) //从对于的group查询缓存值
}
