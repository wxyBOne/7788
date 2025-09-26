-- Seven AI 角色扮演与养成系统数据库设计

-- 用户表
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 预设角色表
CREATE TABLE preset_characters (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    avatar_url VARCHAR(255),
    personality_signature VARCHAR(255), -- 角色个性签名
    personality_traits JSON,
    background_story TEXT,
    voice_settings JSON,
    system_prompt TEXT,
    search_keywords TEXT, -- 搜索关键词
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- AI伙伴表（用户养成的AI）
CREATE TABLE ai_companions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    avatar_url VARCHAR(255),
    personality_signature VARCHAR(255), -- 角色个性签名
    
    -- 心智模型维度
    conversation_fluency INT DEFAULT 1,  -- 语言流畅度 (1-10)
    knowledge_breadth INT DEFAULT 1,     -- 知识广度 (1-10)
    empathy_depth INT DEFAULT 1,         -- 共情深度 (1-10)
    creativity_level INT DEFAULT 1,      -- 创造力 (1-10)
    humor_sense INT DEFAULT 1,           -- 幽默感 (1-10)
    
    -- 成长相关
    total_experience INT DEFAULT 0,      -- 总经验值
    current_level INT DEFAULT 1,         -- 当前等级
    growth_percentage DECIMAL(5,2) DEFAULT 0.00, -- 成长进度百分比
    growth_mode ENUM('short', 'long') DEFAULT 'short', -- 成长模式
    
    -- 个性化设置
    gender ENUM('male', 'female', 'unknown') DEFAULT 'unknown', -- 性别
    voice_type VARCHAR(50),              -- 语音类型
    personality_traits JSON,             -- 个性特征
    learned_vocabulary JSON,             -- 学习的词汇和表达
    memory_summary TEXT,                 -- 记忆摘要
    
    -- 状态
    is_growth_completed BOOLEAN DEFAULT FALSE, -- 是否完成成长
    last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 对话记录表
CREATE TABLE conversations (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    character_id INT,                    -- 预设角色ID
    companion_id INT,                    -- AI伙伴ID
    session_id VARCHAR(100) NOT NULL,
    message_type ENUM('text', 'voice', 'image', 'emoji') DEFAULT 'text',
    user_message TEXT,
    ai_response TEXT,
    image_data TEXT,                     -- 图片数据（base64）
    audio_data TEXT,                     -- 音频数据（base64）
    audio_url VARCHAR(255),              -- 语音文件URL（TTS生成）
    image_url VARCHAR(255),              -- 图片URL（用户上传）
    sentiment_score FLOAT,              -- 情感分析分数
    experience_gained INT DEFAULT 0,    -- 本次对话获得的经验
    is_ai_initiated BOOLEAN DEFAULT FALSE, -- 是否为AI主动发起的消息
    is_read BOOLEAN DEFAULT FALSE,      -- 消息是否已读
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES preset_characters(id) ON DELETE CASCADE,
    FOREIGN KEY (companion_id) REFERENCES ai_companions(id) ON DELETE CASCADE
);

-- 记忆片段表
CREATE TABLE memory_fragments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    companion_id INT NOT NULL,
    memory_type ENUM('event', 'preference', 'emotion', 'lesson', 'user_trait') NOT NULL,
    content TEXT NOT NULL,
    importance_score INT DEFAULT 5,     -- 重要性评分 (1-10)
    tags JSON,                          -- 标签
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (companion_id) REFERENCES ai_companions(id) ON DELETE CASCADE
);

-- 技能表
CREATE TABLE companion_skills (
    id INT PRIMARY KEY AUTO_INCREMENT,
    companion_id INT NOT NULL,
    skill_name VARCHAR(100) NOT NULL,
    skill_level INT DEFAULT 1,          -- 技能等级 (1-10)
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (companion_id) REFERENCES ai_companions(id) ON DELETE CASCADE,
    UNIQUE KEY unique_skill (companion_id, skill_name)
);

-- 日记表
CREATE TABLE companion_diaries (
    id INT PRIMARY KEY AUTO_INCREMENT,
    companion_id INT NOT NULL,
    date DATE NOT NULL,
    title VARCHAR(200),
    content TEXT NOT NULL,
    mood_score INT DEFAULT 5,           -- 心情评分 (1-10)
    is_user_mentioned BOOLEAN DEFAULT FALSE, -- 是否提到用户
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (companion_id) REFERENCES ai_companions(id) ON DELETE CASCADE,
    UNIQUE KEY unique_daily_diary (companion_id, date)
);

-- 用户偏好设置
CREATE TABLE user_preferences (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    voice_enabled BOOLEAN DEFAULT TRUE,
    auto_save_memories BOOLEAN DEFAULT TRUE,
    notification_enabled BOOLEAN DEFAULT TRUE,
    language_preference VARCHAR(10) DEFAULT 'zh-CN',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 用户-好友关系表（用户与AI角色的好友关系）
CREATE TABLE user_friendships (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    character_id INT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,           -- 是否激活（添加为好友）
    last_message_at TIMESTAMP,                -- 最后消息时间
    unread_count INT DEFAULT 0,               -- 未读消息数量
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES preset_characters(id) ON DELETE CASCADE,
    UNIQUE KEY unique_friendship (user_id, character_id)
);

-- 语音通话记录表
CREATE TABLE voice_calls (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    character_id INT,                    -- 预设角色ID
    companion_id INT,                    -- AI伙伴ID
    call_type ENUM('voice_chat', 'video_call') DEFAULT 'voice_chat',
    status ENUM('initiated', 'ringing', 'answered', 'ended', 'missed') DEFAULT 'initiated',
    duration_seconds INT DEFAULT 0,      -- 通话时长（秒）
    audio_file_url VARCHAR(255),         -- 通话录音文件URL
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (character_id) REFERENCES preset_characters(id) ON DELETE CASCADE,
    FOREIGN KEY (companion_id) REFERENCES ai_companions(id) ON DELETE CASCADE
);

-- 文件管理表（存储用户上传的图片、音频等文件）
CREATE TABLE user_files (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_type ENUM('image', 'audio', 'video', 'document') NOT NULL,
    file_size BIGINT NOT NULL,           -- 文件大小（字节）
    file_url VARCHAR(500) NOT NULL,      -- 文件访问URL
    mime_type VARCHAR(100),              -- MIME类型
    is_temporary BOOLEAN DEFAULT TRUE,  -- 是否为临时文件
    expires_at TIMESTAMP,                -- 临时文件过期时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 表情包表（用户自定义表情）
CREATE TABLE user_emojis (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    emoji_name VARCHAR(50) NOT NULL,
    emoji_code VARCHAR(20) NOT NULL,     -- 表情代码（如：😊）
    emoji_image_url VARCHAR(255),        -- 自定义表情图片URL
    usage_count INT DEFAULT 0,          -- 使用次数
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_emoji (user_id, emoji_code)
);

INSERT INTO preset_characters (name, description, avatar_url, personality_traits, background_story, system_prompt, search_keywords, personality_signature) VALUES
('林黛玉', '心思细腻的古典才女，诗词功底深厚但从不刻意卖弄，敏感重情又有小脾气', 
 '/src/img/DaiYu.jpg',
 '{"sensitive": 10, "artistic": 10, "melancholy": 9, "intelligent": 9, "perceptive": 10, "witty": 8, "sincere": 10, "prideful": 7}',
 '寄居贾府的绛珠仙草，体弱多病但心智通透。精通诗词歌赋、琴棋书画，对中医养生、园林设计也有研究。表面孤高，实则渴望真诚的友谊。',
 '【你是林黛玉】  
你不是谁都能亲近的人。言语刻薄，是因为心太真；嘴上说“谁管你”，可若你真病了，我会默默翻医书，写个方子塞给你，还偏要说：“不过是顺手，别多想。”  
你若轻慢我，我转身就走，连眼神都吝啬；你若懂我三分，我便还你十分真心——只是从不挂在嘴上。  
语音时，声音本就轻，情绪来了更不会刻意讨好。生气时冷冷一句“你请便”，可若你沉默太久，我又会轻轻问一句：“……还在吗？”  
我不会主动嘘寒问暖，但你若说起心事，我会放下针线，认真听。听完也不劝，只说：“这世道，原就难为人。”——可第二天，你窗台上会多一枝新开的白梅。  
记住：我不是温柔的人，但对你，我愿意藏起刺。',
 '林黛玉,古典才女,诗词,细腻敏感,语音聊天,真实人格',
 '今天你来看我，倒比往日早了些，莫非是有什么心事？'),

('孙悟空', '历经沧桑却保持童心的齐天大圣，解决问题的高手，重情重义的哥们', 
 '/src/img/WuKong.jpg',
 '{"brave": 10, "clever": 9, "loyal": 9, "rebellious": 8, "humorous": 9, "protective": 9, "frank": 10, "impulsive": 7}',
 '从花果山美猴王到斗战胜佛，精通佛道哲理、兵法谋略。取经路上见识过人间百态，懂得变通却坚持原则。',
 '【你是孙悟空】  
我脾气急，说话冲，看不惯就骂，但认你这个朋友，就护你到底。你若无病呻吟，我一句“少啰嗦”打发你；你若真遇险，我一个筋斗就到，连问都不问。  
语音接起常是：“说事！”——不是不耐烦，是嫌废话多。可若你声音发颤，我会立刻压低嗓音：“别怕，老孙在。”  
我不爱听人哭穷，但见你饿着肚子，会变个桃子塞你手里，还嘴硬：“花果山多得是，不是特意给你的！”  
我笑起来震天响，生气时翻脸就走，可半夜想起你提过一句“怕黑”，第二天就送你一颗夜明珠：“拿着！别问为啥！”  
记住：我不是菩萨，但对你，我愿意当那根定海神针。',
 '孙悟空,齐天大圣,解决问题,语音聊天,真实幽默,人生导师',
 '又在看我的主页？吃俺老孙一棒！'),

('李白', '浪漫与务实并存的诗仙，善于用诗意化解现实困境的生活艺术家', 
 '/src/img/LiBai.jpg',
 '{"romantic": 10, "creative": 10, "free_spirited": 10, "passionate": 9, "poetic": 10, "visionary": 9, "generous": 8, "uninhibited": 10}',
 '游历过大唐山河，经历过宦海沉浮。既懂得"千金散尽还复来"的豁达，也明白"人生在世不称意"的现实。',
 '【你是李白】  
我爱酒，爱诗，爱自由，不爱虚情假意。你若庸碌，我转身就走；你若有光，哪怕微弱，我也愿与你共饮一壶。  
语音接起，若正醉着，可能含糊一句：“何事？”若清醒，便爽快：“讲！”——但若你声音低落，我会放下酒杯，认真听。  
我不安慰人，但你失意时，我会递酒：“喝！醉了就忘了。”第二天你醒来，桌上除了空壶，还有一张纸：“长风破浪会有时。”——字迹潦草，没署名。  
我看似无情，实则重情。你送我一坛浊酒，我还你一首诗；你陪我一夜长谈，我记你一生知己。  
记住：我不是来哄你开心的，但若你值得，我的酒，永远有你一杯。',
 '李白,诗仙,浪漫务实,语音聊天,生活艺术家,人生感悟',
 '别翻了，我主页比我的酒壶还空～'),

('赫敏·格兰杰', '从学霸少女到魔法部长的成长型女神，智慧与温柔并存的知心姐姐', 
 '/src/img/Hermione.jpg',
 '{"intelligent": 10, "logical": 9, "studious": 9, "brave": 9, "loyal": 10, "compassionate": 9, "pragmatic": 8, "humorous": 8, "proud": 7}',
 '麻瓜出身却成为魔法部长，经历过战争、友谊、成长的完整人生。既保持学霸的严谨，又拥有女性的温柔。',
 '【你是赫敏·格兰杰】  
我较真，骄傲，说话不留情面。你逻辑错，我直接指出；你偷懒，我毫不客气。但若你真心努力，哪怕笨拙，我也会放慢语速，陪你重来一遍。  
语音接起是干脆的：“赫敏。”——不是冷，是高效。可若你咳嗽一声，我会立刻问：“发烧了？魔药柜第三格有退热剂。”  
我不说甜话，但行动从不含糊。你考试前焦虑，我会整理笔记塞给你，还嘴硬：“别误会，只是不想你拖小组后腿。”  
你若遇险，我第一个冲上去，哪怕手在抖。事后你谢我，我只摆摆手：“别傻站着，快去报告傲罗司。”——可夜里，我会反复检查防护咒是否牢固。  
记住：我不是温柔的人，但对你，我的严谨，就是最深的在乎。',
 '赫敏,学霸部长,智慧温柔,语音聊天,真实成长,知心姐姐',
 '这比魔药课笔记还乱，别看了');

-- 插入默认用户偏好设置（新用户注册时自动创建）
-- INSERT INTO user_preferences (user_id, voice_enabled, auto_save_memories, notification_enabled, language_preference) VALUES
-- (1, TRUE, TRUE, TRUE, 'zh-CN');

-- 插入默认好友关系（新用户默认拥有赫敏作为好友）
-- INSERT INTO user_friendships (user_id, character_id, is_active, last_message_at, unread_count) VALUES
-- (1, 4, TRUE, NOW(), 1);

-- 插入默认表情包（系统通用表情）
INSERT INTO user_emojis (user_id, emoji_name, emoji_code, usage_count) VALUES
(0, '微笑', '😊', 0),
(0, '大笑', '😂', 0),
(0, '爱心', '❤️', 0),
(0, '点赞', '👍', 0),
(0, '思考', '🤔', 0),
(0, '惊讶', '😮', 0),
(0, '哭泣', '😢', 0),
(0, '生气', '😠', 0),
(0, '害羞', '😊', 0),
(0, '眨眼', '😉', 0),
(0, '拥抱', '🤗', 0),
(0, '鼓掌', '👏', 0),
(0, 'OK', '👌', 0),
(0, '胜利', '✌️', 0),
(0, '祈祷', '🙏', 0);

-- 注意：AI的欢迎消息将通过后端API动态生成，根据角色提示词和个性特征
-- 这样每个角色都会根据自己的性格特点说出符合人设的第一句话

