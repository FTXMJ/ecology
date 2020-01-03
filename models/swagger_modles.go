package models

// Page 分页参数  ---  历史信息
type HostryPageInfo_test struct {
	items___数据列表 []HostryValues_test //数据列表
	page___分页信息  Page_test           //分页信息
}

type Page_test struct {
	total_page__总页数     int //总页数
	current_page___当前页数 int //当前页数
	page_size___每页数据条数  int //每页数据条数
	count___总数据量        int //总数据量
}

type HostryValues_test struct {
	id                      int
	user_id___              string
	current_revenue___本期收入  float64 //上期支出
	current_outlay____本期支出  float64 //本期支出
	opening_balance___上期c余额 float64 //上期c余额
	current_balance___本期余额  float64 //本期余额
	create_date___创建时间      string  //创建时间
	comment___评论_           string  //评论
	tx_id__任务id_            string  //任务id
	account____生态仓库id       int     //生态仓库id
}

type ForceTable_test struct {
	id__id                         int
	level___等级                     string
	low_hold___充值或者升级的低位           int     //低位
	high_hold___高位                 int     //高位
	return_multiple___杠杆           float64 //杠杆
	hold_return_rate____自由算力       float64 //本金自由算力
	recommend_return_rate____直推算力  float64 //直推算力
	team_return_rate____动态算力__团队算力 float64 //动态算力
}

type Ecology_index_ob_test struct {
	usdd___usdd数量                             float64
	ecological_poject___生态项目                  []Formulaindex_test //生态项目
	ecological_poject_bool___是否有生态仓库没有就是false bool
	super_peer___超级节点信息                       SuperPeer_test //超级节点
	super_peer_bool__是否显示超级节点                 bool
}

//页面显示的　超级节点结构
type SuperPeer_test struct {
	usdd___总币数ForceTable_test float64 //总币数ForceTable_test
	level___超级节点的独立属性         string  //超级节点的独立属性
	today_a_bouns___今日分红      float64 // 今日分红
}

// 页面显示的　生态仓库结构
type Formulaindex_test struct {
	id__生态仓库id                   int
	level___等级                   string
	bocked_balance___持币数量        float64 //持币数量
	balance___投资总额               float64 //投资总额
	low_hold___低位                int     //低位
	high_hold___高位               int     //高位
	return_multiple___杠杆         float64 //杠杆
	to_day_rate___今日算力           float64 //今日算力
	hold_return_rate___自由算力      float64 //本金自由算力
	recommend_return_rate___直推算力 float64 //直推算力
	team_return_rate____动态算力     float64 //动态算力
}

// 超级节点算力表
type SuperForceTable_test struct {
	id__id                    int
	level___等级                string
	coin_number_rule___币数     int
	bonus_calculation____分红比例 float64
}

// user coin flow information
type FlowList_test struct {
	Items___数据列表 []Flow_test //数据列表
	Page___分页    Page        //分页信息
}

// user ecology information
type Flow_test struct {
	UserId___用户id              string
	HoldReturnRate___本金自由算力    float64 //本金自由算力
	RecommendReturnRate___直推算力 float64 //直推算力
	TeamReturnRate___动态算力      float64 //动态算力
	Released___已释放             float64 //已释放
	UpdateTime___最后更新时间        string  // 最后更新时间
}

// user ecology information
type UEOBJList_test struct {
	Items___数据列表 []U_E_OBJ_test //数据列表
	Page___分页    Page           //分页信息
}

// user ecology information object
type U_E_OBJ_test struct {
	UserId___用户id              string
	Level___等级                 string
	ReturnMultiple___杠杆        float64 //杠杆
	CoinAll___存币总和             float64 //存币总和
	ToBeReleased___待释放         float64 //待释放
	Released___已释放             float64 //已释放
	HoldReturnRate___本金自由算力    float64 //本金自由算力
	RecommendReturnRate___直推算力 float64 //直推算力
	TeamReturnRate___动态算力      float64 //动态算力
}

// Forces Table
type ForceTable_test_yq struct {
	Id_id                      int     `orm:"column(id);pk;auto" json:"id"`
	Level___等级                 string  `orm:"column(level)" json:"level"`
	LowHold___最低               int     `orm:"column(low_hold)" json:"low_hold"`                           //低位
	HighHold___最高              int     `orm:"column(high_hold)" json:"high_hold"`                         //高位
	ReturnMultiple____杠杆       float64 `orm:"column(return_multiple)" json:"return_multiple"`             //杠杆
	HoldReturnRate____本金自由算力   float64 `orm:"column(hold_return_rate)" json:"hold_return_rate"`           //本金自由算力
	RecommendReturnRate___直推算力 float64 `orm:"column(recommend_return_rate)" json:"recommend_return_rate"` //直推算力
	TeamReturnRate___动态算力      float64 `orm:"column(team_return_rate)" json:"team_return_rate"`           //动态算力
	PictureUrl___图片链接          string  `orm:"column(picture_url)" json:"picture_url"`
}

// user`s peer history table
type PeerHistoryList_test struct {
	Items []PeerHistory_test `json:"items"` //数据列表
	Page  Page               `json:"page"`  //分页信息
}

type PeerHistory_test struct {
	Id___id               int     `json:"id_id"`
	Time时间                string  `json:"time_时间"`
	WholeNetworkTfor全网总算力 float64 `json:"whole_network_tfor_全网总算力"`
	PeerABouns节点总收益       float64 `json:"peer_a_bouns_节点总收益"`
	DiamondsPeer钻石节点数     int     `json:"diamonds_peer_钻石节点数"`
	SuperPeer超级节点数        int     `json:"super_peer_超级节点数"`
	CreationPeer创世节点数     int     `json:"creation_peer_创世节点数"`
}

// peer a_bouns list
type PeerListABouns_test struct {
	Items数据 []PeerAbouns_test `json:"items_数据"` //数据列表
	Page分页  Page              `json:"page_分页"`  //分页信息
}

type PeerAbouns_test struct {
	Id__id       int     `json:"id_id"`
	UserName用户名字 string  `json:"user_name_用户名字"`
	Level等级      string  `json:"level_等级"`
	Tforstfor数量  float64 `json:"tforstfor_数量"`
	Time时间       string  `json:"time_时间"`
}

// user`s peer table
type PeerUserFalse_test struct {
	Items []PeerUser_test `json:"items"` //数据列表
	Page  Page            `json:"page"`  //分页信息
}

type PeerUser_test struct {
	AccountId生态仓库id int     `json:"account_id_生态仓库_id"`
	UserName用户名字    string  `json:"user_name_用户名字"`
	UserId用户id      string  `json:"user_id_用户_id"`
	Level等级         string  `json:"level_等级"`
	State状态         bool    `json:"state_状态"`
	Number数量        float64 `json:"number_数量"`
	UpdateTime更新时间  string  `json:"update_time_更新时间"`
}
