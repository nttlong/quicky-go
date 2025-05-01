package helper_mysql

type HelperMysql struct {
	connectionString string
}

func (m *HelperMysql) GetConnectionString() string {
	return m.connectionString
}
