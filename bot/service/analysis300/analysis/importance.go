package analysis

/*
每个位置数据重要程度，综合评分用
key 对应 HeroDataToName 的序号，value对应其权重(factor)

factor 分为三个档次
0.1: 不重要
0.3: 可以不需要，有更好
0.6: 很重要

除外
场次(total) 特殊处理，综合评分最终会乘一个和 场次置位(total_rank) 对应的权重
factor_total = 0.95 + 0.05 * total_rank

胜率(win) 特殊处理，综合评分最终会乘一个和 胜率置位(win_rank) 对应的权重
factor_win = 0.90 + 0.1 * win_rank
*/

type DataWeight map[int]float64

// HeroTypeWeight 类型权重：肉，近战，辅助，炮台，核心
type HeroTypeWeight [5]float64

// 权重
var (
	a = 0.6
	b = 0.3
	c = 0.1
)

var HeroTypeSlice = []string{"肉", "近战输出", "辅助", "丢丢怪", "红核"}

var ImportanceMap = map[int]string{
	4:  "均每分补刀",
	6:  "均每分k",
	8:  "均每分d",
	10: "均每分a",
	12: "均每分推塔",
	14: "均每分插眼",
	16: "均每分排眼",
	18: "均每分经济",
	21: "均每分输出",
	24: "均每分承伤",
	26: "均每分转换率",
}

// 肉，泛指纯肉，一般作为承伤或者是控制的英雄
var TankImportance = DataWeight{
	4:  c,  // 补刀
	6:  b,  // k
	8:  -b, // d
	10: b,  // a
	12: c,  // 推塔
	14: b,  // 插眼
	16: b,  // 排眼
	18: c,  // 经济
	21: c,  // 输出
	24: a,  // 承伤
	26: c,  // 转换率
}

// 近战输出，拥有大量输出，需要近身的英雄
var ADImportance = DataWeight{
	4:  b,  // 补刀
	6:  a,  // k
	8:  -a, // d
	10: b,  // a
	12: c,  // 推塔
	14: b,  // 插眼
	16: b,  // 排眼
	18: b,  // 经济
	21: a,  // 输出
	24: b,  // 承伤
	26: a,  // 转换率
}

// 辅助，需要出辅助装的英雄
var SupportImportance = DataWeight{
	4:  c,  // 补刀
	6:  c,  // k
	8:  -b, // d
	10: b,  // a
	12: c,  // 推塔
	14: a,  // 插眼
	16: a,  // 排眼
	18: c,  // 经济
	21: c,  // 输出
	24: b,  // 承伤
	26: c,  // 转换率
}

// 炮台，丢丢怪
var PokeImportance = DataWeight{
	4:  b,        // 补刀
	6:  a,        // k
	8:  -a * 1.5, // d
	10: a,        // a
	12: c,        // 推塔
	14: b,        // 插眼
	16: c,        // 排眼
	18: b,        // 经济
	21: a,        // 输出
	24: c,        // 承伤
	26: a,        // 转换率
}

// 核心，泛指红核，需要补刀发育的英雄
var CarryImportance = DataWeight{
	4:  a,        // 补刀
	6:  a,        // k
	8:  -a * 1.5, // d
	10: a,        // a
	12: a,        // 推塔
	14: c,        // 插眼
	16: c,        // 排眼
	18: a,        // 经济
	21: a,        // 输出
	24: c,        // 承伤
	26: a,        // 转换率
}

var importances = [5]DataWeight{TankImportance, ADImportance, SupportImportance, PokeImportance, CarryImportance}

func MergeImportance(weights HeroTypeWeight) (rs DataWeight) {
	rs = map[int]float64{}
	for i, weight := range weights {
		for k, v := range importances[i] {
			rs[k] += v * weight
		}
	}
	return
}

// 各英雄权重类别系数
var HeroFactor = map[string]HeroTypeWeight{
	"雪菜":        {0.9, 0.1, 0, 0, 0},
	"司波达也":      {0, 0, 0, 0, 1},
	"赛贝丝":       {0, 0, 0, 0, 1},
	"猎人娃":       {0, 0, 0, 0, 1},
	"小悟空":       {0, 0.5, 0, 0.5, 0},
	"法师娃":       {0, 0.8, 0, 0.2, 0},
	"伊卡洛斯":      {0, 0, 0, 0, 1},
	"芙兰朵露":      {0, 1, 0, 0, 0},
	"爱蜜莉雅":      {0, 0, 0, 1, 0},
	"小鸡":        {0, 0, 0, 1, 0},
	"本多二代":      {0.3, 0.7, 0, 0, 0},
	"舰队统帅":      {0, 0, 0, 1, 0},
	"魂魄妖梦":      {0.1, 0.9, 0, 0, 0},
	"蒂塔":        {0, 0, 0.9, 0.1, 0},
	"摩尔迦娜":      {0.8, 0.2, 0, 0, 0},
	"一方通行":      {0, 0, 0, 0, 1},
	"迪奥·布兰度":    {0, 1, 0, 0, 0},
	"蝴蝶忍":       {0, 1, 0, 0, 0},
	"战士娃":       {0, 1, 0, 0, 0},
	"佐罗":        {0.1, 0.9, 0, 0, 0},
	"莉娜·因巴斯":    {0, 0.7, 0, 0.3, 0},
	"欧根亲王":      {0, 0, 0, 0, 1},
	"阿库娅":       {0, 0, 0.7, 0.3, 0},
	"妮姆芙":       {0, 0, 1, 0, 0},
	"和真&惠惠":     {0, 0, 0, 1, 0},
	"结城友奈":      {0.5, 0.5, 0, 0, 0},
	"桂木桂马":      {0, 0, 0, 1, 0},
	"美树沙耶加":     {0, 1, 0, 0, 0},
	"天使":        {0, 0, 0, 1, 0},
	"武藤游戏":      {0, 0.3, 0, 0.7, 0},
	"机械神梅普露":    {0.5, 0.5, 0, 0, 0},
	"幻刺露西":      {0, 1, 0, 0, 0},
	"金色之暗":      {0, 1, 0, 0, 0},
	"牧濑红莉栖":     {0, 0, 0, 1, 0},
	"死神":        {0, 0.5, 0, 0.5, 0},
	"姬丝秀忒":      {0.8, 0.2, 0, 0, 0},
	"蓬莱山辉夜":     {0, 0, 0.6, 0.4, 0},
	"帕秋莉·诺蕾姬":   {0, 0, 0, 1, 0},
	"牧师娃":       {0, 0.2, 0.8, 0, 0},
	"伊莎":        {0, 1, 0, 0, 0},
	"射命丸文":      {0, 1, 0, 0, 0},
	"琪露诺":       {0, 1, 0, 0, 0},
	"柏崎星奈":      {0.2, 0.8, 0, 0, 0},
	"逢坂大河":      {0.5, 0.5, 0, 0, 0},
	"蕾姆":        {0, 1, 0, 0, 0},
	"平和岛静雄":     {1, 0, 0, 0, 0},
	"赫拉克勒斯":     {0, 1, 0, 0, 0},
	"肯娘":        {0.5, 0.5, 0, 0, 0},
	"关羽":        {0.5, 0.5, 0, 0, 0},
	"菲特":        {0, 1, 0, 0, 0},
	"我爱罗":       {0, 0.2, 0, 0.8, 0},
	"沢田纲吉":      {0, 1, 0, 0, 0},
	"阿斯托尔福":     {0.7, 0, 0.3, 0, 0},
	"美游":        {0, 0.8, 0.2, 0, 0},
	"貂蝉":        {0, 1, 0, 0, 0},
	"御坂美琴":      {0, 1, 0, 0, 0},
	"环彩羽":       {0, 0, 0, 0, 1},
	"红美铃":       {0.3, 0.7, 0, 0, 0},
	"达克尼斯":      {0.5, 0.5, 0, 0, 0},
	"喔":         {0, 0, 0, 0, 1},
	"少司命":       {0, 1, 0, 0, 0},
	"卡卡西":       {0, 0.5, 0, 0.5, 0},
	"夏娜":        {0, 1, 0, 0, 0},
	"高达":        {0, 0, 0, 0.1, 0.9},
	"楪祈":        {0, 0.5, 0, 0, 0.5},
	"黑贞德":       {0, 1, 0, 0, 0},
	"珂朵莉":       {0, 1, 0, 0, 0},
	"曹焱兵":       {0.2, 0.8, 0, 0, 0},
	"埼玉":        {0.5, 0.5, 0, 0, 0},
	"康娜":        {0, 0, 0, 1, 0},
	"灶门炭治郎":     {0.5, 0.5, 0, 0, 0},
	"七宫智音":      {0, 0, 0, 1, 0},
	"火拳":        {0, 0.3, 0, 0.7, 0},
	"李小狼":       {0, 1, 0, 0, 0},
	"远坂凛":       {0, 0, 0, 1, 0},
	"伊莉雅":       {0, 0, 0, 1, 0},
	"黑":         {0, 1, 0, 0, 0},
	"吉良吉影":      {0, 1, 0, 0, 0},
	"我妻由乃":      {0.3, 0.7, 0, 0, 0},
	"术士娃":       {0, 1, 0, 0, 0},
	"公主":        {0, 1, 0, 0, 0},
	"涅普顿":       {0, 1, 0, 0, 0},
	"栗山未来":      {0.2, 0.8, 0, 0, 0},
	"空":         {0, 0, 1, 0, 0},
	"缇娜":        {0, 0, 0, 1, 0},
	"梅普露":       {0.7, 0, 0.3, 0, 0},
	"菈菈":        {0, 0, 0, 1, 0},
	"赤瞳":        {0, 1, 0, 0, 0},
	"佐助":        {0, 0.3, 0, 0.7, 0},
	"秋山澪":       {0, 0, 1, 0, 0},
	"小鸟游六花":     {0, 1, 0, 0, 0},
	"白岩射手":      {0, 1, 0, 0, 0},
	"伊斯坎达尔":     {0.2, 0.8, 0, 0, 0},
	"神目黑刀":      {0, 1, 0, 0, 0},
	"鸣人":        {0, 1, 0, 0, 0},
	"安兹乌尔恭":     {0, 0.2, 0, 0.8, 0},
	"神裂火织":      {0, 1, 0, 0, 0},
	"夏提雅":       {0.5, 0.5, 0, 0, 0},
	"笠笠笠":       {0, 1, 0, 0, 0},
	"晓美焰":       {0, 0.2, 0, 0.5, 0.3},
	"五更琉璃":      {0, 0, 0, 1, 0},
	"食蜂操祈":      {0, 0, 0.2, 0.8, 0},
	"温蒂":        {0, 0, 1, 0, 0},
	"夜夜":        {0.8, 0.2, 0, 0, 0},
	"阿尔冯斯":      {0.5, 0.5, 0, 0, 0},
	"盗贼娃":       {0, 0.8, 0.2, 0, 0},
	"木更":        {0.2, 0.8, 0, 0, 0},
	"玉藻前":       {0, 0.2, 0.5, 0.3, 0},
	"空条徐伦":      {0, 1, 0, 0, 0},
	"巴依老爷":      {0.5, 0.5, 0, 0, 0},
	"立花响":       {0.2, 0.8, 0, 0, 0},
	"樱满集":       {0.2, 0.8, 0, 0, 0},
	"八云紫":       {0, 0, 0, 1, 0},
	"香风智乃":      {0, 0, 0.3, 0.7, 0},
	"奇犽":        {0, 1, 0, 0, 0},
	"佩姬":        {0, 1, 0, 0, 0},
	"沙耶":        {0, 0, 0, 1, 0},
	"卫宫":        {0.1, 0.3, 0, 0.4, 0.2},
	"亚里亚":       {0, 0.8, 0, 0, 0.2},
	"塞巴斯蒂安":     {0, 1, 0, 0, 0},
	"朝田诗乃":      {0, 0, 0, 1, 0},
	"梦魇":        {0, 0, 0, 0, 1},
	"夜斗":        {0.5, 0.5, 0, 0, 0},
	"真红":        {0, 0.3, 0, 0.7, 0},
	"天草四郎时贞":    {0, 1, 0, 0, 0},
	"白井黑子":      {0, 1, 0, 0, 0},
	"绯斗":        {0.5, 0.5, 0, 0, 0},
	"雅典娜":       {0, 0, 1, 0, 0},
	"黑崎一护":      {0, 1, 0, 0, 0},
	"申屠子夜":      {0, 0.8, 0, 0.2, 0},
	"雾雨魔理沙":     {0, 0, 0, 1, 0},
	"美杜莎":       {0, 0.4, 0.5, 0.1, 0},
	"凯菈":        {0, 1, 0, 0, 0},
	"缇米":        {0, 0, 0, 1, 0},
	"赵云":        {0, 1, 0, 0, 0},
	"君莫笑":       {0, 0.8, 0.2, 0, 0},
	"天道佩恩":      {0.6, 0.2, 0.2, 0, 0},
	"司波深雪":      {0, 0.6, 0.2, 0.2, 0},
	"黑羽快斗":      {0, 0.2, 0, 0.2, 0.6},
	"吉尔伽美什":     {0, 0, 0, 0, 1},
	"黑雪姬":       {0, 0.8, 0.2, 0, 0},
	"西行寺幽幽子":    {0, 1, 0, 0, 0},
	"片翼天使":      {0.6, 0.2, 0.2, 0, 0},
	"贵公子":       {0.6, 0.2, 0.2, 0, 0},
	"奈亚子":       {0, 0, 0, 1, 0},
	"爱丽丝·M":     {0, 0.4, 0, 0.6, 0},
	"冯宝宝":       {0, 1, 0, 0, 0},
	"雅儿贝德":      {0.6, 0.3, 0.1, 0, 0},
	"藤原妹红":      {0, 1, 0, 0, 0},
	"圣大人":       {0, 0, 0, 0, 1},
	"土间埋":       {0, 0.2, 0, 0, 0.8},
	"无头骑士":      {0.8, 0.2, 0, 0, 0},
	"喜羊羊":       {0.5, 0.4, 0.1, 0, 0},
	"白贞":        {0.5, 0, 0.5, 0, 0},
	"亚丝娜":       {0, 1, 0, 0, 0},
	"风音日和":      {0, 0, 1, 0, 0},
	"纳兹":        {0, 1, 0, 0, 0},
	"桐谷和人":      {0, 1, 0, 0, 0},
	"高町奈叶":      {0, 0, 0, 1, 0},
	"炎魔":        {0, 1, 0, 0, 0},
	"美狄亚":       {0, 0.7, 0, 0.3, 0},
	"三千院凪":      {0, 0, 0, 0, 1},
	"琉璃":        {0, 0, 0, 0, 1},
	"莉法":        {0, 1, 0, 0, 0},
	"匹诺曹":       {0, 0, 0, 0, 1},
	"八神疾风":      {0, 0, 0, 1, 0},
	"黑无常":       {0, 1, 0, 0, 0},
	"优吉欧":       {0.5, 0.5, 0, 0, 0},
	"独眼之王":      {0.2, 0.8, 0, 0, 0},
	"大傻":        {0, 0, 0, 0, 1},
	"优克莉伍德":     {0, 0, 1, 0, 0},
	"蕾西亚":       {0, 0, 1, 0, 0},
	"蕾米莉亚":      {0, 0, 0, 1, 0},
	"军姬":        {0, 0, 0, 0, 1},
	"常宣灵":       {0, 0, 0, 1, 0},
	"温蒂·玛贝尔":    {0, 0, 1, 0, 0},
	"尼禄":        {0.2, 0.8, 0, 0, 0},
	"闻人翊悬":      {0.5, 0.5, 0, 0, 0},
	"拿剑爱丽丝":     {0, 0.5, 0, 0, 0},
	"格蕾":        {0, 1, 0, 0, 0},
	"茵蒂克丝":      {0, 0, 0, 0, 1},
	"浅间智":       {0, 0, 0, 0, 1},
	"暴龙兽":       {0, 1, 0, 0, 0},
	"乔鲁诺":       {0, 1, 0, 0, 0},
	"潘多拉":       {0, 0, 0, 1, 0},
	"狱寺隼人":      {0, 0, 0, 1, 0},
	"战场原黑仪":     {0, 1, 0, 0, 0},
	"黑岩射手":      {0, 0, 0, 0, 1},
	"不知火舞":      {0.2, 0.8, 0, 0, 0},
	"鹿目圆香":      {0, 0, 0, 0, 1},
	"诺瓦露":       {0, 1, 0, 0, 0},
	"两仪式":       {0, 1, 0, 0, 0},
	"空条承太郎":     {0, 1, 0, 0, 0},
	"隐居者":       {0, 0, 1, 0, 0},
	"周防尊":       {0.5, 0.5, 0, 0, 0},
	"爱德华":       {0, 1, 0, 0, 0},
	"梦梦":        {0, 0.3, 0, 0.7, 0},
	"剪刀仔":       {0, 1, 0, 0, 0},
	"柒":         {0, 1, 0, 0, 0},
	"王也":        {0, 1, 0, 0, 0},
	"大蛇丸":       {0, 1, 0, 0, 0},
	"白虎":        {0, 0, 0, 1, 0},
	"克子":        {0.2, 0.8, 0, 0, 0},
	"鲁路修":       {0, 0, 0, 0, 1},
	"库丘林":       {0.7, 0.3, 0, 0, 0},
	"天天":        {0, 0, 0, 0, 1},
	"亚瑟王":       {0.5, 0, 0, 0.5, 0},
	"十六夜咲夜":     {0, 0.5, 0, 0.5, 0},
	"托尔":        {0, 0, 0.8, 0.2, 0},
	"富樫勇太":      {0, 0, 0, 1, 0},
	"艾露莎":       {0, 1, 0, 0, 0},
	"路飞":        {0.8, 0.2, 0, 0, 0},
	"骑士娃":       {0.5, 0.5, 0, 0, 0},
	"格雷":        {0, 0.3, 0.7, 0, 0},
	"幻刺莉莉":      {0, 1, 0, 0, 0},
	"纏流子":       {0.2, 0.8, 0, 0, 0},
	"古明地觉&古明地恋": {0, 0, 0, 1, 0},
	"雪斗":        {0, 1, 0, 0, 0},
	"小李":        {0.3, 0.7, 0, 0, 0},
	"白":         {0, 0, 0.2, 0.2, 0.6},
	"波风水门":      {0, 1, 0, 0, 0},
	"风鸣翼":       {0.5, 0.5, 0, 0, 0},
	"乔巴":        {0, 0, 1, 0, 0},
	"绯村剑心":      {0, 1, 0, 0, 0},
	"莉姆":        {0, 1, 0, 0, 0},
	"佐仓杏子":      {0, 1, 0, 0, 0},
	"无名":        {0, 1, 0, 0, 0},
	"立华奏":       {0, 1, 0, 0, 0},
	"涂山红红":      {0.2, 1, 0, 0, 0},
	"梅比斯":       {0, 1, 0, 0, 0},
	"凸守早苗":      {0.5, 0.5, 0, 0, 0},
	"坂田银时":      {0, 1, 0, 0, 0},
	"水银灯":       {0, 0, 0, 0, 1},
	"间桐樱":       {0, 1, 0, 0, 0},
	"小樱":        {0, 0, 0, 1, 0},
	"哈桑":        {0, 1, 0, 0, 0},
	"酷奇":        {0, 1, 0, 0, 0},
	"巴麻美":       {0, 0, 0, 0, 1},
	"迪斯卓尔":      {0, 0, 0, 0, 1},
	"电次":        {0.3, 0.7, 0, 0, 0},
	"神乐":        {0, 1, 0, 0, 0},
	"白胡子":       {0, 1, 0, 0, 0},
	"博丽灵梦":      {0, 0, 0, 1, 0},
	"夜雨声烦":      {0, 1, 0, 0, 0},
	"C.C.":      {0, 0, 1, 0, 0},
	"蛇姬":        {0, 1, 0, 0, 0},
	"歌姬":        {0, 0, 1, 0, 0},
	"苍星石":       {0, 1, 0, 0, 0},
	"石神千空":      {0, 0, 0.8, 0.2, 0},
	"芙蕾雅":       {0.6, 0.2, 0.2, 0, 0},
	"东风谷早苗": {0, 0.2, 0, 0.8, 0},
}
