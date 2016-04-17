namespace cpp tutorial

service PushMsg {
	void notify(1:i32 badge, 2:string clientid, 3:i32 msgtype, 4:string title, 5:string msgcontent)
}
