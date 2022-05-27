package system

import (
	"fmt"
	"time"
)

type SysObject struct {
	// 전체 인스턴스화된 개체를 추적하는 개체 ID
	__GLOBAL_OBJECT_ID int
	__created_time     time.Time
	__object_id        int
	__object_id_other  int
}

func NewSysObject() *SysObject {
	sy := SysObject{}
	sy.__GLOBAL_OBJECT_ID = 0
	sy.__created_time = time.Now()
	sy.__object_id = sy.__GLOBAL_OBJECT_ID
	sy.__GLOBAL_OBJECT_ID = sy.__GLOBAL_OBJECT_ID + 1
	return &sy
}

func (sy SysObject) String() string {
	return fmt.Sprintf("ID:%10d %s", sy.__object_id, sy.__created_time)
}

func Set_req_time(sy *SysObject) {
	return
}

func Get_req_time(sy *SysObject) {
	return
}

func (sy SysObject) __lt__(other SysObject) bool {
	return sy.__object_id < other.__object_id
}

func (sy SysObject) Get_obj_id() int {
	return sy.__object_id
}
