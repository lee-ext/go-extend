package ext

type OptVec[E any] struct {
	opt BytesBitMap
	es  []E
}
