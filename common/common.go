package common

import (
	"encoding/json"
	"fmt"
	"github.com/gavv/httpexpect"
	"gopkg.in/check.v1"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

var headerGlobal = make(map[string]string)
var paramsGlobal = make(map[string]interface{})

func init() {
	ConfigInit()
}

func AddHeaderGlobal(key, value string) {
	headerGlobal[key] = value
}

func AddParamsGlobal(key, value string) {
	paramsGlobal[key] = value
}

// 注册校验
func RegisterCheck(tests ...GoApiTesting) {
	for _, test := range tests {
		var _ = check.Suite(test)
		test.Initialization()
		TestList = append(TestList, test)
	}
}

/**
 * GET接口请求
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func HttpGet(c *check.C, t GoApiTesting) *httpexpect.Response {
	t.Initialization()
	e := httpexpect.New(c, Host)
	uri := urlHandle(t)

	m := t.HandleUrlCode(t)

	// 与全局参数合并
	for k, v := range paramsGlobal {
		if _, ok := m[k]; !ok {
			m[k] = v
		}
	}

	// fmt.Println("----------t.UrlPath() -----" + fmt.Sprint(uri) + "---------------")
	// fmt.Println("----------t.Token() -----" + fmt.Sprint(Token) + "---------------")

	request := e.GET(uri)

	// 参数处理
	for key, value := range m {
		request.WithQuery(key, value)
	}
	// 头处理
	for key, value := range headerGlobal {
		request.WithHeader(key, value)
	}
	r := request.Expect().
		Status(http.StatusOK)
	t.SetResponse(r)

	return r
}

/**
 * POST接口请求
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func HttpPost(c *check.C, t GoApiTesting) *httpexpect.Response {
	t.Initialization()

	contentType := "application/x-www-form-urlencoded;charset=utf-8"
	r := httpPost(c, t, contentType)

	t.SetResponse(r)

	return r
}

func httpPost(c *check.C, t GoApiTesting, contentType string) *httpexpect.Response {
	// 域名
	if Host == "" {
		panic("host is null")
	}
	e := httpexpect.New(c, Host)

	// 请求地址处理
	uri := urlHandle(t)
	// 请求参数处理
	m := t.HandleParam(t)

	// 与全局参数合并
	for k, v := range paramsGlobal {
		if _, ok := m[k]; !ok {
			m[k] = v
		}
	}

	// fmt.Println("-----------uri----" + fmt.Sprint(uri) + "---------------")
	// fmt.Println("------------mmm---" + fmt.Sprint(m) + "---------------")
	// fmt.Println("------------Token---" + fmt.Sprint(Token) + "---------------")
	// fmt.Println("------------contentType---" + fmt.Sprint(contentType) + "---------------")
	request := e.POST(uri)

	// 头处理
	for key, value := range headerGlobal {
		request.WithHeader(key, value)
	}

	r := request.
		WithHeader("ContentType", contentType). // 定义头信息
		WithForm(m).
		Expect().
		Status(http.StatusOK)
	return r
}

func urlHandle(t GoApiTesting) string {
	uri := t.UrlPath()
	if uri == "" {
		panic("uri is null")
	}
	uri = strings.Trim(uri, "/")
	uri = "/" + uri

	return uri
}

/**
 * POST:JSON接口请求
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func HttpPostJson(c *check.C, t GoApiTesting) *httpexpect.Response {
	t.Initialization()

	contentType := "application/json;charset=utf-8"
	r := httpPost(c, t, contentType)

	t.SetResponse(r)

	return r
}

/**
 * 结构体转urlCode
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func structToUrlCode(t GoApiTesting) map[string]interface{} {
	m := make(map[string]interface{})

	sv := reflect.ValueOf(t).Elem()
	st := reflect.TypeOf(t).Elem()
	for i := 0; i < st.NumField(); i++ {
		// fmt.Println("--------key-------" + fmt.Sprint(st.Field(i).Tag.Get("json")) + "---------------")
		// fmt.Println("--------value-------" + fmt.Sprint(sv.Field(i).String()) + "---------------")
		if st.Field(i).Tag.Get("json") != "" {
			m[st.Field(i).Tag.Get("json")] = sv.Field(i).String()
		}
	}

	return m
}

/**
 * 设置TOKEN
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func SetHeaderToken(token string) {
	// fmt.Println("---------token------" + fmt.Sprint(token) + "---------------")
	Token = token
}

/**
 * 从json文件中获取参数
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func ParamByJson(t GoApiTesting, filename string) {
	fileObj, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fileObj.Close()
	content, err := ioutil.ReadAll(fileObj)
	fmt.Println(string(content))
	err = json.Unmarshal(content, t)
	// 随机
	paramRandom(t)
}

/**
 * 全部随机
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func paramRandomByStruct(t GoApiTesting) {
	sv := reflect.ValueOf(t).Elem()
	st := reflect.TypeOf(t).Elem()
	for i := 0; i < st.NumField(); i++ {
		// fmt.Println("--------key-------" + fmt.Sprint(st.Field(i).Tag.Get("json")) + "---------------")
		// fmt.Println("--------value-------" + fmt.Sprint(sv.Field(i).String()) + "---------------")
		// fmt.Println("---------------" + fmt.Sprint(st.Field(i).Type.String()) + "---------------")
		switch st.Field(i).Type.String() {
		case "string":
			str := randomString(8)
			// fmt.Println("---------------" + fmt.Sprint(str) + "---------------")
			sv.Field(i).SetString(str)
		case "int":
			int := rand.Intn(10)
			sv.Field(i).SetInt(int64(int))
		case "int64":
			int := rand.Intn(10)
			sv.Field(i).SetInt(int64(int))
		}
	}

}

func paramRandom(t GoApiTesting) {
	sv := reflect.ValueOf(t).Elem()
	st := reflect.TypeOf(t).Elem()
	for i := 0; i < st.NumField(); i++ {
		// fmt.Println("--------key-------" + fmt.Sprint(st.Field(i).Tag.Get("json")) + "---------------")
		// fmt.Println("--------value-------" + fmt.Sprint(sv.Field(i).String()) + "---------------")
		if sv.Field(i).String() == "auto" {
			// fmt.Println("---------------" + fmt.Sprint(st.Field(i).Type.String()) + "---------------")
			switch st.Field(i).Type.String() {
			case "string":
				str := randomString(8)
				// fmt.Println("---------------" + fmt.Sprint(str) + "---------------")
				sv.Field(i).SetString(str)
			case "int":
				int := rand.Intn(10)
				sv.Field(i).SetInt(int64(int))
			case "int64":
				int := rand.Intn(10)
				sv.Field(i).SetInt(int64(int))
			}
		}
	}
}

func randomString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

/**
 * 处理文件名
 * @Author: cs_shuai
 * @Date: 2020-08-06
 */
func HandleFileName(file string) (string, string) {
	filenameall := path.Base(file)
	filesuffix := path.Ext(file)
	fileprefix := filenameall[0 : len(filenameall)-len(filesuffix)]

	return fileprefix, filesuffix
}

/**
 * 检查json文件类型转换为结构体
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func checkFileTypeToStruct(fileName, path string) *jsonFile {
	fileprefix, filesuffix := HandleFileName(fileName)
	jf := new(jsonFile)
	jf.Path = path + fileName
	jf.FileName = fileName
	jf.FilePrefix = fileprefix
	jf.FileSuffix = filesuffix
	if filesuffix != ".json" {
		panic("文件类型异常: " + fileName)
	}

	return jf
}

/**
 * 通过Json获取测试参数
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func GetParamsByJsonFile(fileName, path string) (result []interface{}) {
	jf := checkFileTypeToStruct(fileName, path)

	filename := RootPath + path + jf.FileName
	fileObj, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fileObj.Close()
	content, err := ioutil.ReadAll(fileObj)
	// fmt.Println(string(content))
	err = json.Unmarshal(content, &result)
	if err != nil {
		panic(err)
	}

	return result
}