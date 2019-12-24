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
	id__id                         int     `orm:"column(id);pk;auto"`
	level___等级                     string  `orm:column(level)`
	low_hold___充值或者升级的低位           int     `orm:column(low_hold)`              //低位
	high_hold___高位                 int     `orm:column(high_hold)`             //高位
	return_multiple___杠杆           float64 `orm:column(return_multiple)`       //杠杆
	hold_return_rate____自由算力       float64 `orm:column(hold_return_rate)`      //本金自由算力
	recommend_return_rate____直推算力  float64 `orm:column(recommend_return_rate)` //直推算力
	team_return_rate____动态算力__团队算力 float64 `orm:column(team_return_rate)`      //动态算力
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
	id__id                    int     `orm:"column(id);pk;auto"`
	level___等级                string  `orm:column(level)`
	coin_number_rule___币数     int     `orm:column(coin_number_rule)`
	bonus_calculation____分红比例 float64 `orm:column(bonus_calculation)`
}
