package sql

type Params struct {
	Driver   string
	Host     *string
	Port     *int
	User     *string
	Password *string
	DB       *string
}
