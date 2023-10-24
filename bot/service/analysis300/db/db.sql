CREATE TABLE
    `matches` (
        `id` int NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
        `match_id` varchar(255) DEFAULT NULL,
        `match_id_int` bigint DEFAULT NULL COMMENT 'int id',
        `m_id` int DEFAULT NULL COMMENT '比赛类型',
        `used_time` int DEFAULT NULL COMMENT '所用时间',
        `create_time` int DEFAULT NULL COMMENT '游戏结束时间戳',
        PRIMARY KEY (`id`),
        UNIQUE KEY `match_id` (`match_id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 0 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '比赛表';

CREATE TABLE
    `players` (
        `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
        `match_id` varchar(255) DEFAULT NULL COMMENT '比赛id',
        `player_id` bigint DEFAULT NULL COMMENT '玩家id',
        `name` varchar(255) DEFAULT NULL COMMENT '玩家昵称',
        `hero_id` int DEFAULT NULL COMMENT '所选英雄id',
        `hero_lv` int DEFAULT NULL COMMENT '结束时英雄等级',
        `side` int DEFAULT NULL COMMENT '1-左下，2-右上',
        `result` int DEFAULT NULL COMMENT '输赢：1-赢，2-输',
        `first_win` int DEFAULT NULL COMMENT '是否首胜：0-否，1；是',
        `summoner_skill1` int DEFAULT NULL COMMENT 'd技能',
        `summoner_skill2` int DEFAULT NULL COMMENT 'f技能',
        `total_money` int DEFAULT NULL COMMENT '结束时经济',
        `kill_unit` int DEFAULT NULL COMMENT '补刀数',
        `used_time` int DEFAULT NULL COMMENT '所用时间',
        `kill_player` int DEFAULT NULL COMMENT '击杀',
        `death` int DEFAULT NULL COMMENT '死亡',
        `assist` int DEFAULT NULL COMMENT '助攻',
        `con_kill_max` int DEFAULT NULL COMMENT '未知',
        `mul_kill_max` int DEFAULT NULL COMMENT '未知',
        `destory_tower` int DEFAULT NULL COMMENT '推塔',
        `treat` int DEFAULT NULL COMMENT '未知',
        `put_eyes` int DEFAULT NULL COMMENT '插眼数',
        `destory_eyes` int DEFAULT NULL COMMENT '排眼数',
        `elo` int DEFAULT NULL COMMENT '团分（已废弃',
        `fv` int DEFAULT NULL COMMENT '竞技力',
        `total_money_side` int DEFAULT NULL COMMENT '己方总经济',
        `total_money_percent` double DEFAULT NULL COMMENT '个人经济占比',
        `make_damage_side` int DEFAULT NULL COMMENT '己方总伤害',
        `make_damage_percent` double DEFAULT NULL COMMENT '个人伤害占比',
        `take_damage_side` int DEFAULT NULL COMMENT '己方总承伤',
        `take_damage_percent` double DEFAULT NULL COMMENT '个人承伤占比',
        `damage_conversion_rate` double DEFAULT NULL COMMENT '伤害转换率',
        `create_time` int DEFAULT NULL COMMENT '游戏结束时间戳',
        PRIMARY KEY (`id`),
        UNIQUE KEY `match_id_2` (`match_id`, `player_id`),
        KEY `hero_id` (`hero_id`),
        KEY `fv` (`fv`),
        KEY `match_id` (`match_id`),
        KEY `player_id` (`player_id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 0 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '玩家战绩表';
