package hertz

const (
	RegTypeConsul string = "consul"
	RegTypeNacos  string = "nacos"
	RegTypeETCD   string = "etcd"
)

const (
	labelMethod     = "method"
	labelStatusCode = "statusCode"
	labelPath       = "path"
	labelBizCode    = "bizStatusCode" // bizStatus

	unknownLabelValue = "unknown"
	succBizStatus     = "0"
)

var defaultBuckets = []float64{
	5000,
	10000,
	25000,
	50000,
	100000,
	250000,
	500000,
	1000000,
	2500000,
	5000000,
	10000000,
	25000000,
	50000000,
	100000000,
}
