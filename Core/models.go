package core


type ModelStruct struct {
}

// func (m *ModelStruct) GetModelObjName() string {
// 	return ""
// }

func (* ModelStruct) isModel() bool {
	return true
}