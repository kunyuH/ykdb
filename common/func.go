package common

//map合并
func MapMergin(mObj ...map[int]string) map[int]string{
	newObj := map[int]string{}
	for _, m := range mObj {
		for k, v := range m {
			newObj[k] = v
		}
	}
	return newObj
}
//map弹出第一个
func MapShift(mObj *map[int]string)(key int,value string) {
	for k, m := range *mObj {
		key = k
		value = m
		break
	}
	delete(*mObj, key)
	return
}
//切片弹出第一个
func SliceShift(sObj *[]string)(sStr string)  {
	newSObj := *sObj
	//取出命令
	sStr = newSObj[:1][0]
	//删除取出的命令
	*sObj = newSObj[1:]
	return
}