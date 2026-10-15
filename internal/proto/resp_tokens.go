package proto

const (
	// RESP payload prefix tokens
	BULKSTRING = '$'
	ARRAY      = '*'
	// RESP	protocol separators
	LF = '\n'
	CR = '\r'

	SERROR = iota * 1
	BULKSTRING_TOK
	ARRAY_TOK
	NULL
	BULKERROR
	VERBATIMSTR
	SET
	PUSHDATA

	// RESP commands for manage Lists, Sets and Vector Sets
	// some of the most used commands for handle LinkedLists
	LPUSH
	RPUSH
	LPUSHX
	LPOP
	RPOP
	BLPOP
	LLEN

	// commands for manager Sets
	SADD
	SREM
	SMEMBERS
	SISMEMBER
	SCARD
	SPOP

	// TODO-> add commands for manipulate Vector Sets

	// standard tokens for identify text and digits
	TEXT
	DIGIT
)

var tokenResolver = map[string]int{
	"lpush":     LPUSH,
	"rpush":     RPUSH,
	"lpushx":    LPUSHX,
	"lpop":      LPOP,
	"rpop":      RPOP,
	"blpop":     BLPOP,
	"llen":      LLEN,
	"sadd":      SADD,
	"srem":      SREM,
	"smembers":  SMEMBERS,
	"sismember": SISMEMBER,
	"scard":     SCARD,
	"spop":      SPOP,
	"*":         ARRAY_TOK,
	"$":         BULKSTRING_TOK,
}
