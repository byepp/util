package regexputil

import "regexp"

// RegexpMatch 用正则表达式匹配第一个()组内的内容
func Match(input string, pattern string) (ret string, matched bool) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}
	mcs := re.FindStringSubmatch(input)
	//ret = re.FindString(input)
	if len(mcs) > 0 {
		matched = true
		ret = mcs[1]
	}
	return
}

// 输入内容是否能够匹配正则表达式
func IsMatch(input string, pattern string) bool {
	matched, err := regexp.MatchString(pattern, input)
	if err != nil {
		return false
	}
	return matched
}
