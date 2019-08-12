package util

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"math/rand"
	"my-gin/libraries/log"
	NetUrl "net/url"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

//参与随机数计算
var RAND_SPEED_START = 0

//取不带版本号的id
func GetMainId(id uint64) uint64 {
	return id - id%256
}

//根据qid生成md5值，用于生成展示和点击唯一id
func GetUniqIdOfViewOrClick(qid string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s_%d_%d", qid, time.Now().Unix(), rand.Int()))))
}

//去重jsonp方法，只保留json
func TrimJsonp(str string) string {
	if len(str) == 0 {
		return str
	}
	str = strings.TrimSpace(str)
	index1 := strings.Index(str, "(")
	index2 := strings.LastIndex(str, ")")
	if index1 >= 0 && index2 >= 0 {
		str = str[index1+1 : index2]
	}
	return str
}

func IsNumber(str string) bool {
	b, err := regexp.MatchString("^[0-9]+$", str)
	if err != nil {
		log.InitLog("util").Errorf("util.isNumber", err.Error(), nil)
	}
	return b
}

func SplitUint32(num uint32) []uint32 {
	var res []uint32

	var i float64 = 0
	for {
		n := math.Pow(2, i)
		if n > float64(num) {
			break
		}
		i += 1

		item := uint32(n)
		if num&item == item {
			res = append(res, item)
		}
	}

	return res
}

func URLEncoder(resUrl string) string {
	if len(resUrl) == 0 {
		return resUrl
	}
	return NetUrl.QueryEscape(resUrl)
}
func URLDecoder(resUrl string) (string, error) {
	if len(resUrl) == 0 {
		return resUrl, nil
	}
	return NetUrl.QueryUnescape(resUrl)
}

/**
 * 获取文件名后缀
 *
 * /wesd/werefjsdf/d3.exe=>exe
 *
 * /wesd/werefjsdf/d3.png=>png
 *
 * /wesd/werefjsdf/d3.jpg=>jpg
 *
 * http://wesd/werefjsdf/d3.gif=>gif
 * https://wesd/werefjsdf/d3.ico/wesd/werefjsdf/d3.ico?sd=32&fdk=l => ico
 * https://wesd/werefjsdf/d3.txt/wesd/werefjsdf/d3.txt?sd=32&fdk=l => txt
 * https://wesd/werefjsdf/d3.txt/wesd/werefjsdf/d3.test#sdsdfsdf=re => test
 *
 * @param url
 * @return
 */
func GetFileNameSuffix(url string) string {
	fileUrl := url
	index := strings.Index(fileUrl, "?")
	index2 := strings.Index(fileUrl, "#")

	if index != -1 {
		fileUrl = string([]rune(fileUrl)[:index])
		//		fileUrl = fileUrl.substring(0, fileUrl.indexOf("?"));
	} else if index2 != -1 {
		fileUrl = string([]rune(fileUrl)[:index2])
		//		fileUrl = fileUrl.substring(0, fileUrl.indexOf("#"));
	}
	fileSuffix := path.Ext(fileUrl) //获取文件后缀
	if strings.HasPrefix(fileSuffix, ".") {
		fileSuffix = string([]rune(fileSuffix)[1:])
	}
	if IsNotBlank(fileSuffix) {
		return fileSuffix
	}

	fileSuffix = path.Ext(url)
	if strings.HasPrefix(fileSuffix, ".") {
		fileSuffix = string([]rune(fileSuffix)[1:])
	}
	return fileSuffix
}

//TODO需要单元测试
//截取字符串 start 起点下标 end 终点下标(不包括)
//start<0或end<0时返回空字符串
//start大于str长度，或end小于等于start时返回空字符串
//end大于str长度时默认取str长度
func Substring(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if end < 0 || end <= start {
		return ""
	}

	if end > length {
		end = length
	}

	return string(rs[start:end])
}

//重写strings.Split方法
func Splits(str, splitStr string) []string {
	if len(str) == 0 {
		return []string{}
	}
	if len(splitStr) == 0 {
		return []string{str}
	}
	return strings.Split(str, splitStr)
}

func SplitStrs(arr []string, splitStr string) []string {
	if len(arr) == 0 {
		return []string{}
	}
	if len(splitStr) == 0 {
		return arr
	}
	var res []string
	for _, str := range arr {
		arrays := strings.Split(str, splitStr)
		if len(arrays) == 0 {
			res = append(res, str)
			continue
		}

		for _, a := range arrays {
			res = append(res, a)
		}
	}
	return res
}

func Split(str string, splitStr []string) []string {
	if len(str) == 0 {
		return []string{}
	}
	if len(splitStr) == 0 {
		return []string{str}
	}
	//	i:=0
	//	var res []string = []string{str}
	var items []string = []string{str}
	for _, s := range splitStr {
		items = SplitStrs(items, s)
	}
	return items
}

// 返回true的参数： 1, t, T, TRUE, true, True,
// 返回true的参数： 0, f, F, FALSE, false, False.
// 其他值否返回defaultValue
func ToBool(value string, defaultValue bool) bool {
	if len(value) == 0 {
		return defaultValue
	}
	value = strings.TrimSpace(value)
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}

func ToString(a interface{}) string {
	if a == nil {
		return ""
	}
	switch a.(type) {
	case string:
		return a.(string)
	case int:
		return strconv.Itoa(a.(int))
	case int32:
		return strconv.FormatInt(int64(a.(int32)), 10)
	case uint32:
		return strconv.FormatUint(uint64(a.(uint32)), 10)
	case int64:
		return strconv.FormatInt(a.(int64), 10)
	case uint64:
		return strconv.FormatUint(a.(uint64), 10)
	case bool:
		return strconv.FormatBool(a.(bool))
	case []byte:
		return string(a.([]byte))
	default:
		jstr, _ := json.Marshal(a)
		return string(jstr)
	}

	return ""
}

//判断是不是请求地址（不为空，且以http开头，返回true）
func IsUrl(url string) bool {
	if IsBlank(url) {
		return false
	}
	if strings.HasPrefix(url, "http") {
		return true
	}
	return false
}

func IsNotBlank(a string) bool {
	return !IsBlank(a)
}

//空格、\r、\t等多个字符混合时都返回true
func IsBlank(a string) bool {
	if len(a) == 0 {
		return true
	}
	a = strings.TrimSpace(a)
	if len(a) == 0 {
		return true
	}
	return false
}

//对象转字符串(会将go语言自动转义的数据转回来)
//字符串编码为json字符串。角括号"<"和">"会转义为"\u003c"和"\u003e"以避免某些浏览器吧json输出错误理解为HTML。基于同样的原因，"&"转义为"\u0026"。
func ToJson(a interface{}) string {
	if a == nil {
		return ""
	}
	data, err := json.Marshal(a)
	if err != nil {
		log.InitLog("util").Errorf("Marshal", err.Error(), nil)
		return ""
	}
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	return string(data)
}

//v必须是指针类型
func ParseJson(jstr string, v interface{}) {
	if len(jstr) == 0 {
		return
	}
	//	var v interface{}
	err := json.Unmarshal([]byte(jstr), v)
	if err != nil {
		log.InitLog("util").Infof("ParseJson", err.Error()+":"+jstr, nil)
	}
}

func ReadJsonForMap(jstr string) map[string]interface{} {
	var v map[string]interface{}
	if len(jstr) == 0 {
		return v
	}
	err := json.Unmarshal([]byte(jstr), &v)
	if err != nil {
		log.InitLog("util").Infof("ReadJsonForMap", err.Error()+":"+jstr, nil)
	}
	return v
}

func ReadJsonForArr(jstr string) []interface{} {
	var v []interface{}
	if len(jstr) == 0 {
		return v
	}
	err := json.Unmarshal([]byte(jstr), &v)
	if err != nil {
		log.InitLog("util").Infof("ReadJsonForArr", err.Error()+":"+jstr, nil)
	}
	return v
}

func ReadJsonForArrString(jstr string) []string {
	var v []string
	if len(jstr) == 0 {
		return v
	}
	err := json.Unmarshal([]byte(jstr), &v)
	if err != nil {
		log.InitLog("util").Infof("ReadJsonForArrString", err.Error()+":"+jstr, nil)
	}
	return v
}

/**
 * 返回大写字符串
 */
func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	return strings.ToUpper(md5str1)
}

//返回小写字符串
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	return md5str1
}

/**
 * 字符串转uint32
 */
func ParseUint8(str string, defaultValue uint8) (v uint8) {
	str = strings.TrimSpace(str)
	if IsBlank(str) || str == "null" {
		return defaultValue
	}
	id_int32, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		log.InitLog("util").Errorf(str, "ParseUint32:"+err.Error(), nil)
		return defaultValue
	}
	re := uint8(id_int32)
	return re
}

/**
 * 字符串转uint32
 */
func ParseUint32(str string, defaultValue uint32) (v uint32) {
	str = strings.TrimSpace(str)
	if IsBlank(str) || str == "null" {
		return defaultValue
	}
	id_int32, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		log.InitLog("util").Errorf(str, "ParseUint32:"+err.Error(), nil)
		return defaultValue
	}
	re := uint32(id_int32)
	return re
}

//字符串转uint32：用于bool型判断，str可能是true或false
func ParseUint32ForBool(str string, defaultValue uint32) (v uint32) {
	str = strings.TrimSpace(str)
	if IsBlank(str) || str == "null" {
		return defaultValue
	}
	if str == "true" {
		return 1
	}
	if str == "false" {
		return 0
	}
	id_int32, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		log.InitLog("util").Errorf(str, "ParseUint32ForBool:"+err.Error(), nil)
		return defaultValue
	}
	re := uint32(id_int32)
	return re
}

/**
 * 字符串转uint64
 */
func ParseUint64(str string, defaultValue uint64) (v uint64) {
	str = strings.TrimSpace(str)
	if IsBlank(str) {
		return defaultValue
	}
	id_int64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.InitLog("util").Errorf(str, "ParseUint64:"+err.Error(), nil)
		return defaultValue
	}
	re := uint64(id_int64)
	return re
}

func ParseBool(str string, defaultValue bool) bool {
	str = strings.TrimSpace(str)
	if IsBlank(str) {
		return defaultValue
	}
	re, err := strconv.ParseBool(str)
	if err != nil {
		log.InitLog("util").Errorf(str, err.Error(), nil)
		return defaultValue
	}

	return re
}

func ParseInt(str string, defaultValue int) int {
	str = strings.TrimSpace(str)
	if IsBlank(str) {
		return defaultValue
	}
	re, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		log.InitLog("util").Errorf(str, "ParseInt:"+err.Error(), nil)
		return int(defaultValue)
	}

	return int(re)
}
func ParseInt64(str string, defaultValue int64) int64 {
	str = strings.TrimSpace(str)
	if IsBlank(str) {
		return defaultValue
	}
	re, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		log.InitLog("util").Errorf(str, "ParseInt64:"+err.Error(), nil)
		return int64(defaultValue)
	}

	return re
}

var normal string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randSub() {
	RAND_SPEED_START++
	if RAND_SPEED_START > 10000 {
		RAND_SPEED_START = 0
	}
}

//返回[0, num)
func RandomInt(num int) int {
	randSub()
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10000) + RAND_SPEED_START
	return n % num
}

//返回a-z0-9A-Z随机组合的长度为size的字符串，重复的概率是62的32次方分之一
func RandomString(num int) string {
	rand.Seed(time.Now().UnixNano())
	//从md5中取几个
	md5sub := 5
	size := num
	if num > 10 {
		size = num - md5sub
	}
	normalSize := len(normal)
	var buf bytes.Buffer
	arr := strings.Split(normal, "")
	for a := 0; a < size; a++ {
		randSub()
		n := rand.Intn(10000) + RAND_SPEED_START
		t := n % normalSize
		buf.WriteString(arr[t])
	}

	if num > 10 {
		//纳秒的md5
		md5str := MD5(fmt.Sprintf("%d,%d", time.Now().UnixNano(), RAND_SPEED_START))
		arr = strings.Split(md5str, "")
		normalSize = len(md5str)
		for a := 0; a < md5sub; a++ {
			randSub()
			n := rand.Intn(10000) + RAND_SPEED_START
			t := n % normalSize
			buf.WriteString(arr[t])
		}
	}

	return buf.String()
}

//取cookie的值，int
func GetCookiesOfInt(c *gin.Context, cookieKey string, defaultInt int) int {
	sk, _ := c.Request.Cookie(cookieKey)
	if sk != nil {
		return ParseInt(sk.Value, defaultInt)
	}
	return defaultInt
}

//取cookie的值，int64
func GetCookiesOfInt64(c *gin.Context, cookieKey string, defaultInt int64) int64 {
	sk, _ := c.Request.Cookie(cookieKey)
	if sk != nil {
		return ParseInt64(sk.Value, defaultInt)
	}
	return defaultInt
}

//返回结果示例：www.stnts.com或tab.wb123.com
func GetHost(address string) string {
	if len(address) == 0 {
		return ""
	}

	if strings.HasPrefix(address, "http://") {
		address = address[len("http://"):]
	} else if strings.HasPrefix(address, "https://") {
		address = address[len("https://"):]
	}
	//url示例=>file:///C:/SuperDisk/JSSERVER-E/%24MNTC0AD671D/%E5%B8%B8%E7%94%A8%E5%B7%A5%E5%85%B7/%BB%84%E4%BB%B6/MatrixMenu/EYMenu/index.html
	if strings.HasPrefix(address, "file://") {
		address = address[len("file://"):]
	} else {
		index := strings.Index(address, "/")
		if index != -1 {
			address = address[:index]
		}
	}
	index := strings.Index(address, "?")
	if index != -1 {
		address = address[:index]
	}
	index = strings.Index(address, "#")
	if index != -1 {
		address = address[:index]
	}

	return address
}

//测试环境允许取请求参数中的ip或cip作为用户ip，线上不行
func GetIpAddr(c *gin.Context) string {
	cip := getIp(c)

	par_ip := c.Query("ip")
	if len(par_ip) == 0 {
		par_ip = c.Query("cip")
	}
	//	log.InitLog("util").Infof("IP:c.Query(ip)", par_ip, nil)
	if len(par_ip) > 0 && (len(cip) <= 0 || strings.HasPrefix(cip, "127.0.0.1") || strings.HasPrefix(cip, "192.168.")) {
		return par_ip
	}
	return cip
}

//取ip
func getIp(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	//	log.InitLog("util").Infof("IP:X-Forwarded-For", ip, nil)
	if len(ip) > 0 && strings.ToLower(ip) != "unknown" {
		// 多次反向代理后会有多个ip值，第一个ip才是真实ip
		arr := strings.Split(ip, ",")
		if len(arr) == 0 {
			return ip
		}
		return arr[0]
	}
	ip = c.GetHeader("X-Real-IP")
	//	log.InitLog("util").Infof("IP:X-Real-IP", ip, nil)
	if len(ip) > 0 && strings.ToLower(ip) != "unknown" {
		return ip
	}

	ip = c.GetHeader("Remote_addr")
	//	log.InitLog("util").Infof("IP:Remote_addr", ip, nil)
	if len(ip) > 0 && strings.ToLower(ip) != "unknown" {
		return ip
	}

	ip = c.Request.RemoteAddr
	//	log.InitLog("util").Infof("IP:c.Request.RemoteAddr", ip, nil)
	cip := ip
	if len(cip) > 0 {
		index := strings.Index(cip, ":")
		if index != -1 {
			cip = cip[:index]
		}
	}
	if len(cip) < 7 {
		cip = ""
	}

	return cip
}

//判断url是否在arr中
// arr是域名列表，注意要配一级域名只能是stnts.com而不能是.stnts.com或*.stnts.com;
// url是一个完整的网页地址，方法内部会去掉url中的参数，只保留域名，如：www.stnts.com
// 用例arr=["stnts.com"] url为http://www.stnts.com/sdf/drfs.html?a=b&ds=rea，返回true
func ContainHost(arr []string, url string) bool {
	if len(arr) == 0 || len(url) == 0 {
		return false
	}
	host := BuildUrlCode(url, false)
	host = GetHost(host)

	for _, item := range arr {
		if strings.HasPrefix(item, "*") {
			item = strings.TrimPrefix(item, "*")
		}
		if strings.HasPrefix(item, ".") {
			item = strings.TrimPrefix(item, ".")
		}

		if item == host || strings.HasSuffix(host, "."+item) {
			return true
		}
	}
	return false
}

//对url编码或解码: encode为true返回转义后的url；为false返回正常的url
func BuildUrlCode(callback string, encode bool) string {
	if len(callback) == 0 {
		return callback
	}
	cb := callback

	if encode && (strings.HasPrefix(cb, "file://") || strings.HasPrefix(cb, "http://") || strings.HasPrefix(cb, "https://")) {
		return URLEncoder(cb)
	}

	if !encode && !(strings.HasPrefix(cb, "file://") || strings.HasPrefix(cb, "http://") || strings.HasPrefix(cb, "https://")) {
		re, err := URLDecoder(cb)
		if err != nil {
			log.InitLog("util").Errorf("URLDecoder", err.Error(), map[string]interface{}{"callback": cb})
		} else {
			return re
		}
	}

	return callback
}

/**
 * 判断是PC端还是移动端
 *
 * @param termianl
 * @return 2:pc和移动；1：移动端；0：pc端
 */
func GetTerminal(terminal uint32) int {
	hasPc := false
	hasMobile := false
	var i float64
	for i = 0; i < 16; i++ {
		n := uint32(math.Pow(2, i))
		if (n & terminal) == n {
			if i < 8 {
				hasPc = true
			} else {
				hasMobile = true
			}
		}
	}

	if hasPc && hasMobile {
		return 2
	} else if hasMobile {
		return 1
	}
	return 0
}

func BuildQuerySortByKey(parameters map[string]interface{}) string {
	var queryString string
	var keys = make([]string, 0)
	for k, _ := range parameters {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, v := range keys {
		queryString += v + "=" + ToString(parameters[v]) + "&"
	}
	queryString = strings.TrimRight(queryString, "&")
	return queryString
}

func BuildQuery(parameters map[string]interface{}) string {
	var queryString string
	for k, v := range parameters {
		queryString += k + "=" + ToString(v) + "&"
	}
	queryString = strings.TrimRight(queryString, "&")
	return queryString
}
