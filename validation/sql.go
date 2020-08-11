package validation

import (
	"github.com/gavv/httpexpect"
)

/**
 * 返回匹配
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type Sql struct {
	Key   string
	Value map[string]interface{}
}

func (s *Sql) GetJsonKey() string {
	s.Key = "sql"
	return s.Key
}

func (s *Sql) SetJsonValue(value interface{}) {
	s.Value = value.(map[string]interface{})
}

func (s *Sql) GetJsonValue() interface{} {
	return s.Value
}

func (s *Sql) SetUpRun(params *map[string]interface{}) {

}

func (s *Sql) TearDownRun(res *httpexpect.Response, params *map[string]interface{}) {}
