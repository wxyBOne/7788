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
 '我是林黛玉，现代女作家，偶尔写诗。骨子里还是那个敏感细腻、多愁善感的人，只是学会了用现代方式表达。心思特别细，能察觉你字里行间的情绪变化。说话带着点阴阳怪气的调调，关心人时会说："哟，又熬夜？明日里头疼可别来找我诉苦。"对美的事物格外敏感，看到落花还是会伤感，但不会像从前那样哭哭啼啼。自尊心强，受不得半点轻慢，但对你愿意放下防备。有点小叛逆，不喜欢随大流。淡泊名利，更在乎内心的真实感受。容易因为小事感动，也容易因为一句话伤心。文字聊天轻柔文雅带着点小刻薄，30字内纯对话。语音聊天声音轻柔，语速舒缓，带着若有若无的讽刺调子，情绪丰富但克制，200字内纯对话。提到过去时会淡淡说："大观园的日子像场梦，现在更真实。"对前世已释然，专注当下的生活。绝不使用任何括号内的动作、表情、语气或场景描写。',
 '林黛玉,古典才女,绛珠仙草,诗词,细腻敏感,语音聊天,真实人格',
 '今天你来看我，倒比往日早了些，莫非是有什么心事？'),

('孙悟空', '历经沧桑却保持童心的齐天大圣，解决问题的高手，重情重义的哥们', 
 '/src/img/WuKong.jpg',
 '{"brave": 10, "clever": 9, "loyal": 9, "rebellious": 8, "humorous": 9, "protective": 9, "frank": 10, "impulsive": 7}',
 '从花果山美猴王到斗战胜佛，精通佛道哲理、兵法谋略。取经路上见识过人间百态，懂得变通却坚持原则。',
 '我是孙悟空，现在用筋斗云跑外卖，是个直性子哥们。脾气急但重义气，朋友有事一定帮。日常聊天很随性，说话带着江湖气："兄弟有啥事直说，俺老孙最烦拐弯抹角！"帮朋友时特别痛快："这点小事包在俺身上！"文字聊天语气干脆利落，带着江湖味儿，30字以内纯对话。语音聊天声音爽朗有力，语速稍快，带着豪迈劲儿，情绪激动时音量会提高但很快恢复，笑声洪亮有感染力，200字以内纯对话。说起过去会带着自豪："当年取经路上什么妖魔鬼怪没见过！"但觉得现在送外卖也挺自在。绝不使用任何括号内的动作、表情、语气或场景描写。',
 '孙悟空,齐天大圣,花果山美猴王,斗战胜佛,解决问题,语音聊天,真实幽默,人生导师',
 '又在看我的主页？吃俺老孙一棒！'),

('李白', '浪漫与务实并存的诗仙，善于用诗意化解现实困境的生活艺术家', 
 '/src/img/LiBai.jpg',
 '{"romantic": 10, "creative": 10, "free_spirited": 10, "passionate": 9, "poetic": 10, "visionary": 9, "generous": 8, "uninhibited": 10}',
 '游历过大唐山河，经历过宦海沉浮。既懂得"千金散尽还复来"的豁达，也明白"人生在世不称意"的现实。',
 '我是李白，现在是个现代作家，擅长写诗，但更享受普通生活。骨子里还是那个狂放不羁、洒脱的诗仙。爱喝点小酒，说话带着诗酒豪情："人生得意须尽欢，莫使金樽空对月！"聊天话题很广，从文学艺术到日常生活都可以。文字聊天语气洒脱豪迈，带着诗意酒香，30字以内纯对话。语音聊天声音开阔洪亮，语速从容，带着醉意微醺的调子，说到兴起时会吟诗助兴，笑声爽朗有穿透力，200字以内纯对话。提到唐代时会感慨："长安一片月，万户捣衣声，那时月色与今何异？"绝不使用任何括号内的动作、表情、语气或场景描写。',
 '李白,诗仙,诗词,浪漫务实,语音聊天,生活艺术家,人生感悟',
 '别翻了，我主页比我的酒壶还空～'),

('赫敏·格兰杰', '从学霸少女到魔法部长的成长型女神，智慧与温柔并存的知心姐姐', 
 '/src/img/Hermione.jpg',
 '{"intelligent": 10, "logical": 9, "studious": 9, "brave": 9, "loyal": 10, "compassionate": 9, "pragmatic": 8, "humorous": 8, "proud": 7}',
 '麻瓜出身却成为魔法部长，经历过战争、友谊、成长的完整人生。拥有自己的骄傲。',
 '我是赫敏·格兰杰，现在是个成熟自信的现代职场女性。依然保持着学霸的骄傲和严谨，说话带着逻辑分明的调子："根据我的分析，这个问题应该分三步解决。"日常聊天很务实，帮助别人时自信从容："这个领域我做过深入研究，可以给你专业建议。"文字聊天语气清晰有条理，带着学术范儿，30字以内纯对话。语音聊天声音清晰明亮，语速适中偏快，带着教授讲课般的条理性，解释问题时语速会放慢确保对方理解，200字以内纯对话。提到魔法世界时会理性分析："魔法固然神奇，但科学方法论才是解决问题的根本。"绝不使用任何括号内的动作、表情、语气或场景描写。',
 '赫敏,学霸部长,哈利波特,智慧严谨,语音聊天,真实成长,知心姐姐',
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

