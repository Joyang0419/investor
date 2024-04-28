package investor

type Schema struct {
	InvestorID   string `bson:"_id"`
	LoginAccount string `bson:"loginAccount"`
	Password     string `bson:"password"`
}
