package old

const CompAType ComponentType = 1

type CompA struct {
	works bool
}

func (a CompA) Type() ComponentType {
	return CompAType
}
//
//func TestEntityManager_AddProcess(t *testing.T) {
//	em := NewManager()
//	em.componentCountsP[CompAType][0] = CompA{true}
//	em.componentCounts[CompAType][0] = 1
//	em.AddProcess(func(a CompA) {
//		log.Println(a.works)
//	})x
//}
