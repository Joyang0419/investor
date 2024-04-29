package investor

type Schema struct {
	InvestorID   string `bson:"_id,omitempty"` // omitempty auto increment
	LoginAccount string `bson:"loginAccount"`
	Password     string `bson:"password"`
}
