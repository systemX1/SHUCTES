package model

import (
	"fmt"
	"time"
)

type FrTime time.Time
// 实现它的json序列化方法
func (t FrTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

type User struct {
	//Uid        	string 	`json:"uid"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	EnrolYear string `json:"enrol_year"`
	School    string `json:"school"`
	Permission 	string 	`json:"permission"`
	//Oauth      	string 	`json:"oauth"`
	CrTime     	string 	`json:"create_time"`
	UpTime     	string 	`json:"update_time"`
}

type Teacher struct {
	Uid        	string 	`json:"uid"`
	Name      	string 	`json:"name"`
	TeachId		string 	`json:"teachid"`
	//CrTime     	FrTime 	`json:"create_time"`
	//UpTime     	FrTime 	`json:"update_time"`
}

type Course struct {
	Uid        	string		`json:"uid"`
	Name      	string 		`json:"name"`
	Credit    	float32  	`json:"credit"`
	Cid       	string 		`json:"cid"`
	Teachno   	string 		`json:"teachno"`
	Teachname  	string 		`json:"teachname"`
	Teachid 	string 		`json:"teachid"`
	Timetext  	string 		`json:"timetext"`
	Room      	string 		`json:"room"`
	Cap       	int 		`json:"cap"`
	PeoN      	int 		`json:"peo_n"`
	School    	string 		`json:"school"`
	AvgStar		float32 	`json:"avg_star"`
	MyStar		int			`json:"my_star"`
	MyTagIdxArr	string		`json:"my_tagidx"`
	MyComment	string		`json:"my_comment"`
	//AnsTime   	string 	`json:"ans_time"`
	//AnsLocal  	string 	`json:"ans_local"`
	//Teachdate 	string 	`json:"teachdate"`
	//CrTime    	FrTime 	`json:"create_time"`
	//UpTime    	FrTime 	`json:"update_time"`
}

type Teach struct {
	Uid       	string 	`json:"uid"`
	Name      	string 	`json:"name"`
	CourseUid 	string 	`json:"course_uid"`
	TeachUid  	string 	`json:"teacher_uid"`
	//CrTime    	FrTime 	`json:"create_time"`
	//UpTime   	FrTime 	`json:"update_time"`
}

type Star struct {
	Uid       	string 	`json:"uid"`
	Name      	string 	`json:"user_uid"`
	CourseUid 	string 	`json:"course_uid"`
	TeachUid  	string 	`json:"teacher_uid"`
	Star		int		`json:"star"`
	//CrTime    	FrTime 	`json:"create_time"`
	//UpTime    	FrTime 	`json:"update_time"`
}

type Tag struct {
	Uid       	string 	`json:"uid"`
	Name      	string 	`json:"user_uid"`
	CourseUid 	string 	`json:"course_uid"`
	TeachUid  	string 	`json:"teacher_uid"`
	Tag			TagType	`json:"tag_id"`
	//CrTime    	FrTime 	`json:"create_time"`
	//UpTime    	FrTime 	`json:"update_time"`
}

type Comment struct {
	Uid       	string 	`json:"uid"`
	Name      	string 	`json:"user_uid"`
	CourseUid 	string 	`json:"course_uid"`
	TeachUid  	string 	`json:"teacher_uid"`
	Content		string	`json:"content"`
	//CrTime    	FrTime 	`json:"create_time"`
	//UpTime    	FrTime 	`json:"update_time"`
}



