package mongox

type Config struct {
	// mongodb://username:password@example1.com,example2.com,example3.com/?replicaSet=test&w=majority&wtimeoutMS=5000
	MongoUri string `mapstructure:"mongo_uri" json:"mongo_uri" toml:"mongo_uri"`
}
