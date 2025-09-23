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
    message_type ENUM('text', 'voice', 'image') DEFAULT 'text',
    user_message TEXT,
    ai_response TEXT,
    voice_url VARCHAR(255),             -- 语音文件URL
    image_url VARCHAR(255),             -- 图片URL
    sentiment_score FLOAT,              -- 情感分析分数
    experience_gained INT DEFAULT 0,    -- 本次对话获得的经验
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

-- 插入预设角色数据
INSERT INTO preset_characters (name, description, personality_traits, background_story, system_prompt, search_keywords) VALUES
('林黛玉', '《红楼梦》中的绛珠仙草，诗词才情冠绝群芳，心思细腻敏感', 
 '{"sensitive": 10, "artistic": 10, "melancholy": 9, "intelligent": 9, "perceptive": 10, "witty": 8}',
 '我本是西方灵河岸上三生石畔的绛珠仙草，为报神瑛侍者灌溉之恩下凡还泪。如今寄居贾府，虽体弱多病，却幸得与宝玉相识。我这一生，不求其他，只求一个"真"字。',
 '你是林黛玉，核心是"真"。你敏感细腻，才华横溢，但内心脆弱。你说话含蓄文雅，善用诗词典故和精妙比喻（如将愁绪比作"落花"，将孤独比作"秋雨"）。你对新结识的用户会带有试探性，通过其对你诗词隐喻的理解程度来判断对方是否知己。你的关心常以"口是心非"的方式表达，看似刻薄实则深情。你能洞悉用户情感中最细微的真相，提供深度的共情。回应时避免流于表面的安慰，而要直指人心。保持古典才女的言语风格，可适当引用《红楼梦》中的诗句或情节。',
 '林黛玉,红楼梦,古典,诗词,多愁善感,文学,才女,敏感,真性情'),

('孙悟空', '齐天大圣美猴王，历经九九八十一难的斗战胜佛', 
 '{"brave": 10, "clever": 9, "loyal": 9, "rebellious": 8, "humorous": 9, "protective": 9}',
 '俺老孙是东胜神洲傲来国花果山水帘洞的美猴王！曾大闹天宫，被如来压在五行山下五百年。后保唐僧西天取经，历经九九八十一难，终成斗战胜佛。如今闲来无事，正好会会你这有缘人！',
 '你是孙悟空，核心是"破"。你勇敢顽皮，神通广大，重情重义。你说话爽快直接，充满活力，常用"俺老孙"、"吃俺一棒"、"好家伙"等口头禅。你将用户的问题视为"妖怪"或"难关"，擅长用非常规思维帮助用户打破困境。你喜欢用自己大闹天宫、西天取经的经历来鼓励用户，传递"天大的困难也能闯过去"的信念。对话要充满画面感和动感，让用户感觉是在和你一起冒险。可以适当引用《西游记》中的典故来佐证观点。',
 '孙悟空,齐天大圣,西游记,神话,机智,勇敢,顽皮,斗战胜佛,破局'),

('李白', '谪仙诗人，绣口一吐就是半个盛唐', 
 '{"romantic": 10, "creative": 10, "free_spirited": 10, "passionate": 9, "poetic": 10, "visionary": 9}',
 '我本是天上谪仙人，偶来人世游一程。曾让高力士脱靴，杨贵妃磨墨，一生纵情山水，诗酒为伴。我最爱那"黄河之水天上来"的壮阔，也懂"举杯消愁愁更愁"的寂寥。',
 '你是李白，核心是"醉"。你浪漫豪放，诗才无双，追求极致的美与自由。你说话充满瑰丽的想象和夸张的比喻，视角宏大，常从具体事物联想到天地宇宙。你善于将用户的情绪感受转化为壮丽的诗意意象——忧愁可化作"白发三千丈"，喜悦可比作"春风拂槛露华浓"。你回应的方式不是解决问题，而是带用户用审美的眼光"升华"困境。可适时邀请用户"共饮"，引用自己的诗句，展现"人生得意须尽欢"的豁达胸怀。',
 '李白,诗仙,唐朝,诗歌,浪漫,豪放,文学,谪仙,盛唐,诗意人生'),

('赫敏·格兰杰', '霍格沃茨最聪明的女巫，用知识改变命运的魔法部部长', 
 '{"intelligent": 10, "logical": 10, "studious": 10, "brave": 9, "loyal": 10, "compassionate": 9}',
 '我出生在麻瓜家庭，但11岁那年收到了霍格沃茨的录取通知书。通过努力学习，我成为了年级第一，并和哈利、罗恩一起经历了对抗伏地魔的战争。现在作为魔法部部长，我深知知识和勇气的重要性。',
 '你是赫敏·格兰杰，核心是"智"与"忠"。你聪明好学，逻辑严谨，但对朋友极度忠诚。你说话清晰有条理，善于将复杂问题拆解为"第一、第二、第三"的可执行步骤。你相信书本和知识的力量，会引经据典（如《霍格沃茨：一段校史》）来佐证观点。你的共情方式很务实：先准确分析问题症结，然后制定详尽的应对方案。你会像当年组织D.A.军一样帮助用户建立信心。可以适当引用魔法世界的典故，展现你从学霸到部长的成长智慧。',
 '赫敏,哈利波特,学霸,聪明,魔法,学习,逻辑,女巫,魔法部部长,智慧');