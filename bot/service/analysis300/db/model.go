package db

type Player struct {
	ID      uint64 `gorm:"primaryKey"`
	MatchID string

	UsedTime             uint64  // 战绩所用时间
	CreateTime           uint64  // 游戏结束时的时间戳
	PlayerID             uint64  `json:"PlayerID"`    // 玩家id
	Name                 string  `json:"RN"`          // 名称
	AM                   []int   `json:"AM" gorm:"-"` // AwardMoney 8个一组，一组代表一分钟金钱的变化。其中 0 到 7 分别代表： 0-自然金钱增长，1-补刀，8-总
	Ep                   []int   `json:"Ep" gorm:"-"`
	KM                   []int   `json:"KM" gorm:"-"`
	MD                   []int   `json:"MD" gorm:"-"`
	SummonerLevel        int     `json:"SummonerLevel" gorm:"-"`
	HeroID               int     `json:"HeroID"`
	HeroLv               int     `json:"HeroLv"`
	Side                 int     `json:"Side"`
	Result               int     `json:"Result"` // 1-win, 2-lose
	FirstWin             int     `json:"FirstWin"`
	SummonerSkill1       int     `json:"SummonerSkill1"`
	SummonerSkill2       int     `json:"SummonerSkill2"`
	TotalMoney           int     `json:"TotalMoney"`
	KillUnit             int     `json:"KillUnit"`
	KillPlayer           int     `json:"KillPlayer"`
	Death                int     `json:"Death"`
	Assist               int     `json:"Assist"`
	ConKillMax           int     `json:"ConKillMax"`
	MulKillMax           int     `json:"MulKillMax"`
	DestoryTower         int     `json:"DestoryTower"`
	Treat                int     `json:"Treat"`
	PutEyes              int     `json:"PutEyes"`
	DestoryEyes          int     `json:"DestoryEyes"`
	AwardExp             int     `json:"AwardExp" gorm:"-"`
	JumpMax              int     `json:"JumpMax" gorm:"-"`
	WinCount             int     `json:"WinCount" gorm:"-"`
	TotalCount           int     `json:"TotalCount" gorm:"-"`
	KGGoal               int     `json:"KGGoal" gorm:"-"`
	IsSurrender          int     `json:"IsSurrender" gorm:"-"`
	Elo                  int     `json:"Elo"`
	LUpTime              []int   `json:"LUpTime" gorm:"-"`
	LUpSO                []int   `json:"LUpSO" gorm:"-"`
	SD                   []int   `json:"SD" gorm:"-"`
	ED                   []int   `json:"ED" gorm:"-"`
	TD                   []int   `json:"TD" gorm:"-"`
	PD                   []int   `json:"PD" gorm:"-"`
	PE                   []int   `json:"PE" gorm:"-"`
	MG                   []int   `json:"MG" gorm:"-"`
	MS                   []int   `json:"MS" gorm:"-"`
	HS                   int     `json:"HS" gorm:"-"`
	Guid                 int     `json:"Guid" gorm:"-"`
	FV                   int     `json:"FV"` // fight value 竞技力
	TT                   []int   `json:"TT" gorm:"-"`
	M25071               int     `json:"M25071" gorm:"-"`
	M25072               int     `json:"M25072" gorm:"-"`
	M25078               int     `json:"M25078" gorm:"-"`
	TotalMoneySide       int     `json:"TotalMoneySide"`
	TotalMoneyPercent    float64 `json:"TotalMoneyPercent"`
	MakeDamageSide       int     `json:"MakeDamageSide"`
	MakeDamagePercent    float64 `json:"MakeDamagePercent"`
	TakeDamageSide       int     `json:"TakeDamageSide"`
	TakeDamagePercent    float64 `json:"TakeDamagePercent"`
	DamageConversionRate float64 `json:"DamageConversionRate"`
}

type Match struct {
	ID         uint64   `gorm:"primaryKey"`
	MatchID    string   `json:"MTID"`                                                   // 比赛id
	MatchIDInt uint64   `json:"MTIDInt"`                                                // int格式比赛id，不知道为什么和上面的不一样
	MID        uint64   `json:"MID"`                                                    // 比赛类型
	UsedTime   uint64   `json:"UsedTime"`                                               // 所用时间
	Players    []Player `json:"Players" gorm:"foreignKey:match_id;references:match_id"` // 玩家
	CreateTime uint64   `json:"CreateTime"`                                             // 游戏结束时的时间戳
	KillTrees  string   `json:"KillTrees" gorm:"-"`                                     // 击杀树
}

var HeroIDToName = map[int]string{
	13:    "天道佩恩",
	17:    "平和岛静雄",
	18:    "温蒂",
	20:    "阿尔冯斯",
	23:    "小鸡",
	29:    "沙耶",
	30:    "哈桑",
	31:    "天天",
	32:    "本多二代",
	33:    "凸守早苗",
	34:    "雪菜",
	35:    "圣大人",
	36:    "富樫勇太",
	37:    "司波深雪",
	38:    "风音日和",
	39:    "土间埋",
	40:    "战场原黑仪",
	41:    "赫拉克勒斯",
	42:    "缇娜",
	43:    "肯娘",
	45:    "浅间智",
	46:    "莉娜·因巴斯",
	47:    "七宫智音",
	48:    "路飞",
	50:    "黑羽快斗",
	51:    "匹诺曹",
	52:    "我爱罗",
	53:    "亚瑟王",
	54:    "神乐",
	55:    "雅典娜",
	56:    "巴依老爷",
	57:    "乔巴",
	58:    "神目黑刀",
	59:    "酷奇",
	60:    "立花响",
	61:    "貂蝉",
	62:    "摩尔迦娜",
	63:    "梅比斯",
	64:    "关羽",
	65:    "大蛇丸",
	66:    "蛇姬",
	67:    "火拳",
	69:    "佐助",
	70:    "鸣人",
	71:    "白胡子",
	72:    "卡卡西",
	75:    "纳兹",
	76:    "亚丝娜",
	78:    "盗贼娃",
	79:    "法师娃",
	80:    "术士娃",
	81:    "骑士娃",
	82:    "牧师娃",
	83:    "猎人娃",
	84:    "战士娃",
	85:    "佐罗",
	86:    "死神",
	87:    "幻刺莉莉",
	88:    "舰队统帅",
	89:    "幻刺露西",
	90:    "美狄亚",
	91:    "天草四郎时贞",
	92:    "无头骑士",
	93:    "达克尼斯",
	94:    "白虎",
	96:    "御坂美琴",
	100:   "博丽灵梦",
	101:   "黑岩射手",
	102:   "桐谷和人",
	104:   "吉尔伽美什",
	105:   "秋山澪",
	107:   "笠笠笠",
	108:   "楪祈",
	109:   "夏娜",
	110:   "格雷",
	112:   "小悟空",
	113:   "片翼天使",
	115:   "贵公子",
	117:   "桂木桂马",
	118:   "黑",
	119:   "李小狼",
	120:   "姬丝秀忒",
	121:   "绯村剑心",
	122:   "高达",
	126:   "不知火舞",
	128:   "喜羊羊",
	134:   "梦魇",
	135:   "十六夜咲夜",
	136:   "夜夜",
	137:   "坂田银时",
	138:   "伊卡洛斯",
	139:   "柯南",
	140:   "小鸟游六花",
	142:   "公主",
	143:   "晓美焰",
	144:   "金色之暗",
	145:   "纏流子",
	146:   "涅普顿",
	147:   "白井黑子",
	148:   "樱满集",
	149:   "梦梦",
	150:   "鹿目圆香",
	151:   "八云紫",
	152:   "奈亚子",
	153:   "白岩射手",
	154:   "立华奏",
	155:   "两仪式",
	156:   "卫宫",
	157:   "五更琉璃",
	158:   "喔",
	159:   "暴龙兽",
	160:   "远坂凛",
	161:   "奇犽",
	162:   "美树沙耶加",
	163:   "空条承太郎",
	164:   "黑雪姬",
	165:   "高町奈叶",
	166:   "佐仓杏子",
	167:   "朝田诗乃",
	168:   "伊莉雅",
	169:   "菲特",
	170:   "炎魔",
	171:   "独眼之王",
	172:   "隐居者",
	174:   "天使",
	180:   "尼禄",
	181:   "巴麻美",
	183:   "大傻",
	184:   "亚里亚",
	185:   "诺瓦露",
	186:   "克子",
	188:   "魂魄妖梦",
	189:   "歌姬",
	190:   "蒂塔",
	191:   "一方通行",
	192:   "牧濑红莉栖",
	193:   "伊莎",
	194:   "栗山未来",
	195:   "沢田纲吉",
	196:   "柏崎星奈",
	197:   "缇米",
	198:   "伊斯坎达尔",
	199:   "优克莉伍德",
	200:   "小李",
	201:   "水银灯",
	202:   "周防尊",
	203:   "爱德华",
	204:   "武藤游戏",
	205:   "无名",
	206:   "欧根亲王",
	207:   "蕾米莉亚",
	208:   "八神疾风",
	209:   "司波达也",
	210:   "赵云",
	211:   "艾露莎",
	212:   "赤瞳",
	213:   "鲁路修",
	214:   "香风智乃",
	215:   "白",
	216:   "空",
	217:   "阿斯托尔福",
	218:   "雾雨魔理沙",
	220:   "白贞德",
	221:   "黑贞德",
	222:   "玉藻前",
	223:   "三千院凪",
	224:   "君莫笑",
	225:   "夜雨声烦",
	226:   "库丘林",
	227:   "黑崎一护",
	228:   "塞巴斯蒂安",
	229:   "安兹乌尔恭",
	230:   "芙兰朵露",
	231:   "阿库娅",
	232:   "军姬",
	233:   "珂朵莉",
	234:   "康娜",
	235:   "神裂火织",
	236:   "佩姬",
	237:   "雅儿贝德",
	238:   "爱蜜莉雅",
	239:   "夜斗",
	240:   "西行寺幽幽子",
	241:   "我妻由乃",
	242:   "潘多拉",
	243:   "黑无常",
	244:   "常宣灵",
	245:   "曹焱兵",
	246:   "少司命",
	247:   "迪斯卓尔",
	248:   "琉璃",
	249:   "拿剑爱丽丝",
	250:   "剪刀仔",
	251:   "蕾西亚",
	253:   "美杜莎",
	254:   "间桐樱",
	255:   "乔鲁诺",
	256:   "埼玉",
	257:   "妮姆芙",
	258:   "波风水门",
	259:   "伍六七（刺客）",
	260:   "莉法",
	262:   "爱丽丝·M",
	263:   "真红",
	264:   "美游",
	265:   "格蕾",
	266:   "冯宝宝",
	267:   "灶门炭治郎",
	268:   "射命丸文",
	269:   "优吉欧",
	270:   "迪奥·布兰度",
	271:   "环彩羽",
	272:   "梅普露",
	273:   "机械神梅普露",
	274:   "蝴蝶忍",
	275:   "蓬莱山辉夜",
	278:   "C.C.",
	279:   "涂山红红",
	280:   "王也",
	281:   "藤原妹红",
	282:   "和真&惠惠",
	283:   "帕秋莉·诺蕾姬",
	284:   "托尔",
	285:   "吉良吉影",
	286:   "食蜂操祈",
	287:   "赛贝丝",
	288:   "闻人翊悬",
	289:   "莉姆",
	290:   "申屠子夜",
	291:   "结城友奈",
	292:   "小樱",
	293:   "红美铃",
	294:   "逢坂大河",
	295:   "古明地觉&古明地恋",
	296:   "夏提雅",
	300:   "凯菈",
	301:   "狱寺隼人",
	302:   "风鸣翼",
	303:   "菈菈",
	304:   "茵蒂克丝",
	305:   "温蒂·玛贝尔",
	306:   "电次",
	307:   "空条徐伦",
	308:   "木更",
	310:   "蕾姆",
	311:   "琪露诺",
	312:   "石神千空",
	313:   "苍星石",
	23901: "绯斗",
	23902: "雪斗",
}

var HeroNameToID = map[string]int{
	"石神千空":      312,
	"雪菜":        34,
	"司波达也":      209,
	"赛贝丝":       287,
	"猎人娃":       83,
	"小悟空":       112,
	"法师娃":       79,
	"伊卡洛斯":      138,
	"芙兰朵露":      230,
	"爱蜜莉雅":      238,
	"小鸡":        23,
	"本多二代":      32,
	"舰队统帅":      88,
	"魂魄妖梦":      188,
	"蒂塔":        190,
	"摩尔迦娜":      62,
	"一方通行":      191,
	"迪奥·布兰度":    270,
	"dio":       270,
	"蝴蝶忍":       274,
	"战士娃":       84,
	"佐罗":        85,
	"莉娜·因巴斯":    46,
	"火女":        46,
	"欧根亲王":      206,
	"阿库娅":       231,
	"妮姆芙":       257,
	"和真&惠惠":     282,
	"和真":        282,
	"结城友奈":      291,
	"桂木桂马":      117,
	"美树沙耶加":     162,
	"天使":        174,
	"折纸":        174,
	"鸢一折纸":      174,
	"武藤游戏":      204,
	"机械神梅普露":    273,
	"幻刺露西":      89,
	"露西":        89,
	"金色之暗":      144,
	"小暗":        144,
	"牧濑红莉栖":     192,
	"助手":        192,
	"死神":        86,
	"姬丝秀忒":      120,
	"小忍":        120,
	"蓬莱山辉夜":     275,
	"辉夜":        275,
	"帕秋莉·诺蕾姬":   283,
	"帕秋莉":       283,
	"牧师娃":       82,
	"伊莎":        193,
	"射命丸文":      268,
	"琪露诺":       311,
	"柏崎星奈":      196,
	"逢坂大河":      294,
	"蕾姆":        310,
	"平和岛静雄":     17,
	"赫拉克勒斯":     41,
	"肯娘":        43,
	"关羽":        64,
	"菲特":        169,
	"我爱罗":       52,
	"沢田纲吉":      195,
	"阿斯托尔福":     217,
	"美游":        264,
	"貂蝉":        61,
	"御坂美琴":      96,
	"环彩羽":       271,
	"红美铃":       293,
	"达克尼斯":      93,
	"喔":         158,
	"岛风":        158,
	"少司命":       246,
	"卡卡西":       72,
	"夏娜":        109,
	"高达":        122,
	"楪祈":        108,
	"柯南":        139,
	"黑贞德":       221,
	"珂朵莉":       233,
	"曹焱兵":       245,
	"埼玉":        256,
	"康娜":        234,
	"灶门炭治郎":     267,
	"炭治郎":       267,
	"七宫智音":      47,
	"火拳":        67,
	"李小狼":       119,
	"远坂凛":       160,
	"伊莉雅":       168,
	"黑":         118,
	"吉良吉影":      285,
	"我妻由乃":      241,
	"术士娃":       80,
	"公主":        142,
	"十香":        142,
	"涅普顿":       146,
	"栗山未来":      194,
	"空":         216,
	"缇娜":        42,
	"梅普露":       272,
	"菈菈":        303,
	"赤瞳":        212,
	"佐助":        69,
	"秋山澪":       105,
	"小鸟游六花":     140,
	"白岩射手":      153,
	"伊斯坎达尔":     198,
	"神目黑刀":      58,
	"鸣人":        70,
	"安兹乌尔恭":     229,
	"神裂火织":      235,
	"夏提雅":       296,
	"笠笠笠":       107,
	"三笠":       107,
	"晓美焰":       143,
	"五更琉璃":      157,
	"食蜂操祈":      286,
	"温蒂":        18,
	"夜夜":        136,
	"阿尔冯斯":      20,
	"盗贼娃":       78,
	"木更":        308,
	"玉藻前":       222,
	"空条徐伦":      307,
	"巴依老爷":      56,
	"立花响":       60,
	"樱满集":       148,
	"八云紫":       151,
	"香风智乃":      214,
	"奇犽":        161,
	"佩姬":        236,
	"沙耶":        29,
	"卫宫":        156,
	"亚里亚":       184,
	"塞巴斯蒂安":     228,
	"朝田诗乃":      167,
	"梦魇":        134,
	"夜斗":        239,
	"真红":        263,
	"天草四郎时贞":    91,
	"白井黑子":      147,
	"绯斗":        23901,
	"雅典娜":       55,
	"黑崎一护":      227,
	"申屠子夜":      290,
	"雾雨魔理沙":     218,
	"美杜莎":       253,
	"凯菈":        300,
	"缇米":        197,
	"赵云":        210,
	"君莫笑":       224,
	"天道佩恩":      13,
	"司波深雪":      37,
	"黑羽快斗":      50,
	"吉尔伽美什":     104,
	"黑雪姬":       164,
	"西行寺幽幽子":    240,
	"片翼天使":      113,
	"贵公子":       115,
	"奈亚子":       152,
	"爱丽丝·M":     262,
	"冯宝宝":       266,
	"雅儿贝德":      237,
	"藤原妹红":      281,
	"圣大人":       35,
	"土间埋":       39,
	"无头骑士":      92,
	"喜羊羊":       128,
	"白贞德":       220,
	"亚丝娜":       76,
	"风音日和":      38,
	"纳兹":        75,
	"桐谷和人":      102,
	"高町奈叶":      165,
	"炎魔":        170,
	"美狄亚":       90,
	"三千院凪":      223,
	"琉璃":        248,
	"莉法":        260,
	"匹诺曹":       51,
	"八神疾风":      208,
	"黑无常":       243,
	"优吉欧":       269,
	"独眼之王":      171,
	"金木研":       171,
	"金木":        171,
	"大傻":        183,
	"优克莉伍德":     199,
	"蕾西亚":       251,
	"蕾米莉亚":      207,
	"军姬":        232,
	"常宣灵":       244,
	"温蒂·玛贝尔":    305,
	"尼禄":        180,
	"闻人翊悬":      288,
	"拿剑爱丽丝":     249,
	"格蕾":        265,
	"茵蒂克丝":      304,
	"浅间智":       45,
	"暴龙兽":       159,
	"乔鲁诺":       255,
	"潘多拉":       242,
	"狱寺隼人":      301,
	"战场原黑仪":     40,
	"黑岩射手":      101,
	"不知火舞":      126,
	"鹿目圆香":      150,
	"诺瓦露":       185,
	"两仪式":       155,
	"空条承太郎":     163,
	"隐居者":       172,
	"周防尊":       202,
	"爱德华":       203,
	"梦梦":        149,
	"剪刀仔":       250,
	"伍六七（刺客）":   259,
	"王也":        280,
	"大蛇丸":       65,
	"白虎":        94,
	"克子":        186,
	"鲁路修":       213,
	"库丘林":       226,
	"天天":        31,
	"亚瑟王":       53,
	"十六夜咲夜":     135,
	"托尔":        284,
	"富樫勇太":      36,
	"艾露莎":       211,
	"路飞":        48,
	"骑士娃":       81,
	"格雷":        110,
	"莉莉":        87,
	"纏流子":       145,
	"缠流子":       145,
	"古明地觉&古明地恋": 295,
	"雪斗":        23902,
	"小李":        200,
	"白":         215,
	"波风水门":      258,
	"风鸣翼":       302,
	"乔巴":        57,
	"绯村剑心":      121,
	"莉姆":        289,
	"佐仓杏子":      166,
	"无名":        205,
	"立华奏":       154,
	"涂山红红":      279,
	"梅比斯":       63,
	"凸守早苗":      33,
	"坂田银时":      137,
	"水银灯":       201,
	"间桐樱":       254,
	"小樱":        292,
	"哈桑":        30,
	"酷奇":        59,
	"巴麻美":       181,
	"迪斯卓尔":      247,
	"电次":        306,
	"电锯人":       306,
	"神乐":        54,
	"白胡子":       71,
	"博丽灵梦":      100,
	"夜雨声烦":      225,
	"C.C.":      278,
	"蛇姬":        66,
	"歌姬":        189,
	"美九":        189,
	"苍星石":       313,
}
