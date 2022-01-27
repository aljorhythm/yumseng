package rooms

type UserInfoHeap []*UserInfo

func (h UserInfoHeap) Len() int {
	return len(h)
}

func (h UserInfoHeap) Less(i, j int) bool {
	return h[i].Points < h[j].Points
}

func (h UserInfoHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *UserInfoHeap) Push(x interface{}) {
	*h = append(*h, x.(*UserInfo))
}

func (h *UserInfoHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
