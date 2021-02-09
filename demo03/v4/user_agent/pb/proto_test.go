package pb

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

type Message struct {
	MessageMeta string  // 元数据
	MessageName string  // 消息名
	TypeMap map[string]string  //proto类型 | golang类型
}

type Method struct {
	MethodName string // 方法名
	Param string // 参数
	Returns string // 返回值
}

func TestNewUserClient(t *testing.T) {
	var messages []Message
	f, err := os.OpenFile("user.proto", os.O_RDONLY,0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	contentByte,err :=ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	contentStr := string(contentByte)

	//获取Message
	messages = MatchMessage(contentStr)
	for _, message := range messages {
		fmt.Printf("messageName : %s %+v \n", message.MessageName, message.TypeMap)
	}
}

func MatchMessage(contentStr string) (messages []Message) {
	ret := regexp.MustCompile(`message[\s]*([\S]*)[\s\S]*?}`)
	result := ret.FindAllStringSubmatch(contentStr, -1)
	messages = make([]Message, len(result))
	for i, str := range result {
		messages[i].MessageMeta = str[0]
		messages[i].MessageName = str[1]
		//[\s]*([\S]*)[\s]*([\S]*)[\s]*=[\s]*[\d]*
		subRet := regexp.MustCompile(`(.*)[\s]*=[\s]*[\d]*`)
		subResult := subRet.FindAllStringSubmatch(str[0],-1)
		messages[i].TypeMap = make(map[string]string, len(subResult))
		for _, str2 := range subResult {
			dataInfo := strings.Split(strings.TrimSpace(str2[1]), " ")
			if dataInfo[0] != "repeated" {
				messages[i].TypeMap[dataInfo[1]] = dataInfo[0]
			} else {
				messages[i].TypeMap[dataInfo[2]] = "[]"+dataInfo[1]
			}
		}
	}
	return messages
}

func TestRegExp(t *testing.T) {
	//str := "abc a7c mfc cat aMc azc cba"
	//// 解析、编译正则表达式
	////ret := regexp.MustCompile(`a.c`)  	// `` : 表示使用原生字符串
	//ret := regexp.MustCompile(`a[0-9a-z]c`)  	// `` : 表示使用原生字符串
	//
	//// 提取需要信息
	//alls := ret.FindAllStringSubmatch(str, -1)
	//result := ret.FindAllString(str, -1)
	//fmt.Println("alls:", alls)
	//fmt.Printf("result: %v \n", result)

	flysnowRegexp := regexp.MustCompile(`^http://www.flysnow.org/([\d]{4})/([\d]{2})/([\d]{2})/([\w-]+).html$`)
	params := flysnowRegexp.FindStringSubmatch("http://www.flysnow.org/2018/01/20/golang-goquery-examples-selector.html")

	for _,param :=range params {
		fmt.Println(param)
	}
}