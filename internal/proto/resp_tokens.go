package proto

const (
	// RESP payload prefix tokens
	SSTRING    = '+'
	INTEGER    = ':'
	BULKSTRING = '$'
	ARRAY      = '*'
	// RESP	protocol separators
	LF = '\n'

	SERROR = iota * 1
	SSTRING_TOK
	INTEGER_TOK
	BULKSTRING_TOK
	ARRAY_TOK
	NULL
	BOOLEAN
	FLOAT
	BIGNUMBER
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
	"LPUSH":     LPUSH,
	"RPUSH":     RPUSH,
	"LPUSHX":    LPUSHX,
	"LPOP":      LPOP,
	"RPOP":      RPOP,
	"BLPOP":     BLPOP,
	"LLEN":      LLEN,
	"SADD":      SADD,
	"SREM":      SREM,
	"SMEMBERS":  SMEMBERS,
	"SISMEMBER": SISMEMBER,
	"SCARD":     SCARD,
	"SPOP":      SPOP,
	"*":         ARRAY_TOK,
	"+":         SSTRING_TOK,
	"$":         BULKSTRING_TOK,
	":":         INTEGER_TOK,
}
