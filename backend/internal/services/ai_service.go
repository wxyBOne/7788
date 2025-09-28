// Package services 提供AI相关的业务逻辑服务
package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// AIService AI服务，处理LLM对话、语音识别和语音合成
type AIService struct {
	apiKey  string       // API密钥
	baseURL string       // API基础URL
	model   string       // 使用的模型名称
	client  *http.Client // HTTP客户端
}

// ChatRequest LLM对话请求结构
type ChatRequest struct {
	Model       string    `json:"model"`       // 模型名称
	Messages    []Message `json:"messages"`    // 消息列表
	Temperature float64   `json:"temperature"` // 温度参数
	MaxTokens   int       `json:"max_tokens"`  // 最大token数
}

// Message 单条消息结构
type Message struct {
	Role    string `json:"role"`    // 角色：user/assistant/system
	Content string `json:"content"` // 消息内容
}

// ChatResponse LLM对话响应结构
type ChatResponse struct {
	Choices []Choice `json:"choices"` // 选择列表
	Usage   Usage    `json:"usage"`   // 使用统计
}

// Choice 对话选择结构
type Choice struct {
	Message Message `json:"message"` // 消息内容
	Index   int     `json:"index"`   // 索引
}

// Usage token使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`     // 提示token数
	CompletionTokens int `json:"completion_tokens"` // 完成token数
	TotalTokens      int `json:"total_tokens"`      // 总token数
}

// ASRRequest 语音识别请求结构
type ASRRequest struct {
	Model string `json:"model"` // 模型名称
	Audio struct {
		Format string `json:"format"` // 音频格式
		URL    string `json:"url"`    // 音频URL
	} `json:"audio"`
	Language string `json:"language,omitempty"` // 语言参数
}

// ASRResponse 语音识别响应结构
type ASRResponse struct {
	Reqid     string `json:"reqid"`     // 请求ID
	Operation string `json:"operation"` // 操作类型
	Data      struct {
		Result struct {
			Text string `json:"text"` // 识别结果文本
		} `json:"result"`
	} `json:"data"`
}

// TTSRequest 语音合成请求结构
type TTSRequest struct {
	Audio struct {
		VoiceType  string  `json:"voice_type"`  // 音色类型
		Encoding   string  `json:"encoding"`    // 编码格式
		SpeedRatio float64 `json:"speed_ratio"` // 语速比例
	} `json:"audio"`
	Request struct {
		Text      string `json:"text"`                 // 合成文本
		MaxTokens int    `json:"max_tokens,omitempty"` // 最大token数
	} `json:"request"`
}

// TTSResponse 语音合成响应结构
type TTSResponse struct {
	Reqid     string `json:"reqid"`     // 请求ID
	Operation string `json:"operation"` // 操作类型
	Sequence  int    `json:"sequence"`  // 序列号
	Data      string `json:"data"`      // base64编码的音频数据
}

type VisionRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type VisionResponse struct {
	Choices []VisionChoice `json:"choices"`
}

type VisionChoice struct {
	Message VisionMessage `json:"message"`
}

type VisionMessage struct {
	Content string `json:"content"`
}

// NewAIService 创建AI服务实例
func NewAIService(apiKey, baseURL, model string) *AIService {
	return &AIService{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 10 * time.Second, // HTTP客户端超时时间
		},
	}
}

// ChatWithLLM 与LLM进行对话
func (s *AIService) ChatWithLLM(messages []Message, model string, temperature float64, messageType string) (string, error) {
	// 处理表情消息
	messages = s.processEmojiMessages(messages)

	// 添加消息类型标识
	if messageType == "voice" {
		// 在系统消息中添加语音通话标识
		systemMessage := Message{
			Role:    "system",
			Content: "这是一次语音通话，请用符合角色的口吻以及语气自然地和用户对话，请口语化而不是书面语。绝不使用任何括号内的动作、表情、语气或场景描写。请保持回复简洁，控制在60字以内。",
		}
		messages = append([]Message{systemMessage}, messages...)
	}

	// 暂时返回模拟响应，等API配置完成后替换
	if s.apiKey == "" {
		return "你好！我是AI助手，很高兴为您服务。", nil
	}

	// 使用配置的模型，如果传入的model为空则使用默认模型
	usedModel := model
	if usedModel == "" {
		usedModel = s.model
	}

	req := ChatRequest{
		Model:       usedModel,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   40, // 减少到40，让AI生成更简洁的回复，适合语音通话
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", s.baseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("AI API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from AI")
}

// AnalyzeImage 分析图片
func (s *AIService) AnalyzeImage(imageBase64 string, prompt string) (string, error) {
	// 暂时返回模拟响应
	return "我看到了一张图片，但暂时无法详细分析。", nil
}

// SpeechToText 语音转文字 (ASR)
func (s *AIService) SpeechToText(audioData []byte) (string, error) {
	if len(audioData) == 0 || s.apiKey == "" {
		return "", nil
	}

	start := time.Now()

	// 上传音频文件到公网URL
	publicURL, err := s.uploadToPublicURL(audioData)
	if err != nil {
		return "", err
	}
	fmt.Printf("上传耗时: %v\n", time.Since(start))

	// 调用真实的ASR API
	asrStart := time.Now()
	result, err := s.callQiniuASRAPI(publicURL)
	if err != nil {
		return "", err
	}
	fmt.Printf("ASR耗时: %v\n", time.Since(asrStart))
	fmt.Printf("总耗时: %v\n", time.Since(start))

	return result, nil
}

// callQiniuASRAPI 调用七牛云ASR API（快速失败策略）
func (s *AIService) callQiniuASRAPI(audioURL string) (string, error) {
	// 只尝试MP3格式，快速失败
	fmt.Printf("尝试使用 mp3 格式进行ASR识别...\n")

	// 构建ASR请求
	req := ASRRequest{
		Model: "asr",
		Audio: struct {
			Format string `json:"format"`
			URL    string `json:"url"`
		}{
			Format: "mp3",
			URL:    audioURL,
		},
		Language: "zh-CN", // 指定中文语言
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("ASR请求构建失败: %v", err)
	}

	// 尝试ASR请求（内部已有重试机制）
	result, err := s.tryASRRequest(s.baseURL+"/voice/asr", reqBody)
	if err == nil {
		fmt.Printf("ASR识别成功: %s\n", result)
		return result, nil
	}

	// 如果MP3失败，快速尝试WAV格式
	fmt.Printf("MP3格式失败，快速尝试WAV格式...\n")
	req.Audio.Format = "wav"
	reqBody, err = json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("WAV格式请求构建失败: %v", err)
	}

	result, err = s.tryASRRequest(s.baseURL+"/voice/asr", reqBody)
	if err == nil {
		fmt.Printf("WAV格式ASR识别成功: %s\n", result)
		return result, nil
	}

	return "", fmt.Errorf("所有格式的ASR请求都失败了")
}

// tryASRRequest 尝试ASR请求（带重试机制）
func (s *AIService) tryASRRequest(url string, reqBody []byte) (string, error) {
	maxRetries := 2               // 减少到2次重试
	retryDelay := 1 * time.Second // 减少延迟到1秒

	for attempt := 1; attempt <= maxRetries; attempt++ {
		start := time.Now()

		// 创建带超时的HTTP客户端
		client := &http.Client{
			Timeout: 8 * time.Second, // 8秒超时，平衡响应速度和成功率
		}

		// 发送HTTP请求
		httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		if err != nil {
			if attempt == maxRetries {
				return "", fmt.Errorf("failed to create ASR request: %w", err)
			}
			fmt.Printf("ASR请求创建失败 (尝试 %d/%d): %v\n", attempt, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

		resp, err := client.Do(httpReq)
		if err != nil {
			if attempt == maxRetries {
				return "", fmt.Errorf("ASR HTTP request failed: %w", err)
			}
			fmt.Printf("ASR请求失败 (尝试 %d/%d): %v\n", attempt, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}
		defer resp.Body.Close()

		fmt.Printf("ASR请求耗时: %v (URL: %s, 尝试 %d/%d)\n", time.Since(start), url, attempt, maxRetries)

		if resp.StatusCode != 200 {
			respBody, _ := io.ReadAll(resp.Body)
			if attempt == maxRetries {
				return "", fmt.Errorf("ASR API returned status %d: %s", resp.StatusCode, string(respBody))
			}
			fmt.Printf("ASR API返回错误状态 %d (尝试 %d/%d): %s\n", resp.StatusCode, attempt, maxRetries, string(respBody))
			time.Sleep(retryDelay)
			continue
		}

		// 解析响应
		var asrResp ASRResponse
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			if attempt == maxRetries {
				return "", fmt.Errorf("failed to read ASR response: %w", err)
			}
			fmt.Printf("ASR响应读取失败 (尝试 %d/%d): %v\n", attempt, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		if err := json.Unmarshal(respBody, &asrResp); err != nil {
			if attempt == maxRetries {
				return "", fmt.Errorf("failed to parse ASR response: %w", err)
			}
			fmt.Printf("ASR响应解析失败 (尝试 %d/%d): %v\n", attempt, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		// 提取识别文本
		if asrResp.Data.Result.Text != "" {
			fmt.Printf("ASR识别成功: %s\n", asrResp.Data.Result.Text)
			return asrResp.Data.Result.Text, nil
		}

		// 如果没有文本但请求成功，返回空字符串而不是错误
		fmt.Printf("ASR识别结果为空，但请求成功\n")
		return "", nil
	}

	return "", fmt.Errorf("ASR请求在 %d 次尝试后仍然失败", maxRetries)
}

// uploadToPublicURL 上传音频到公网URL
func (s *AIService) uploadToPublicURL(audioData []byte) (string, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("上传耗时: %v (文件大小: %d bytes)\n", time.Since(start), len(audioData))
	}()

	// 检查音频长度，太短的音频可能识别不准确
	if len(audioData) < 8000 { // 少于0.5秒的音频
		return "", fmt.Errorf("音频太短 (%d bytes)，可能影响识别准确性", len(audioData))
	}

	// 如果文件太大，尝试压缩
	if len(audioData) > 100000 { // 100KB
		fmt.Printf("文件较大 (%d bytes)，尝试优化上传\n", len(audioData))
	}

	// 将PCM数据转换为WAV格式
	wavData, err := s.convertPCMToWAV(audioData)
	if err != nil {
		return "", fmt.Errorf("PCM转WAV失败: %w", err)
	}

	// 创建multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件字段 - 使用WAV格式
	filename := fmt.Sprintf("audio_%d.wav", time.Now().Unix())

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = part.Write(wavData)
	if err != nil {
		return "", fmt.Errorf("failed to write file data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 尝试多个上传服务
	uploadServices := []string{
		"https://0x0.st",
		"https://file.io",
	}

	var lastErr error
	for _, serviceURL := range uploadServices {
		fmt.Printf("尝试上传到: %s\n", serviceURL)

		// 重新创建请求（因为body已经被读取）
		var newBuf bytes.Buffer
		newWriter := multipart.NewWriter(&newBuf)

		newPart, err := newWriter.CreateFormFile("file", filename)
		if err != nil {
			lastErr = fmt.Errorf("failed to create form file: %w", err)
			continue
		}

		_, err = newPart.Write(wavData)
		if err != nil {
			lastErr = fmt.Errorf("failed to write file data: %w", err)
			continue
		}

		err = newWriter.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to close writer: %w", err)
			continue
		}

		req, err := http.NewRequest("POST", serviceURL, &newBuf)
		if err != nil {
			lastErr = fmt.Errorf("failed to create upload request: %w", err)
			continue
		}

		req.Header.Set("Content-Type", newWriter.FormDataContentType())

		// 为上传创建专门的客户端，合理超时时间
		uploadClient := &http.Client{
			Timeout: 5 * time.Second, // 减少上传超时时间，提高响应速度
		}

		resp, err := uploadClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("upload request failed: %w", err)
			continue
		}

		if resp.StatusCode == 200 {
			// 读取响应
			respBody, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				lastErr = fmt.Errorf("failed to read upload response: %w", err)
				continue
			}

			// 处理不同服务的响应格式
			var publicURL string
			if serviceURL == "https://file.io" {
				// file.io返回JSON格式
				var result struct {
					Success bool   `json:"success"`
					Key     string `json:"key"`
					Link    string `json:"link"`
				}
				if err := json.Unmarshal(respBody, &result); err == nil && result.Success {
					publicURL = result.Link
				}
			} else {
				// 0x0.st返回直接的URL
				publicURL = strings.TrimSpace(string(respBody))
			}

			// 验证URL格式
			if strings.HasPrefix(publicURL, "http") {
				fmt.Printf("上传成功到: %s, URL: %s\n", serviceURL, publicURL)
				return publicURL, nil
			}
		}

		if resp.Body != nil {
			resp.Body.Close()
		}
		lastErr = fmt.Errorf("upload failed with status %d", resp.StatusCode)
	}

	return "", fmt.Errorf("所有上传服务都失败，最后错误: %w", lastErr)
}

// SpeechToTextWithCharacter 带角色人设的语音转文字
func (s *AIService) SpeechToTextWithCharacter(audioData []byte, characterName string) (string, error) {
	// 调用ASR
	text, err := s.SpeechToText(audioData)
	if err != nil {
		return "", nil // 返回空字符串，让后端处理
	}

	// 如果识别结果为空或太短，返回空字符串
	if strings.TrimSpace(text) == "" || len(strings.TrimSpace(text)) < 2 {
		return "", nil
	}

	return text, nil
}

// getNoiseResponseForCharacter 根据角色人设返回噪音响应
func (s *AIService) getNoiseResponseForCharacter(characterName string) string {
	switch characterName {
	case "林黛玉":
		return "哎呀，这声音怎么听不清呢？你那边是不是太吵了？"
	case "孙悟空":
		return "兄弟，你这声音俺老孙听不清啊，是不是环境太吵了？"
	case "李白":
		return "这声音如雾里看花，听不真切，莫非是环境嘈杂？"
	case "赫敏·格兰杰":
		return "抱歉，环境噪音太大，我听不清楚你说什么。"
	default:
		return "抱歉，我听不清楚你说什么，可能是环境太吵了。"
	}
}

// GetAPIKey 获取API密钥
func (s *AIService) GetAPIKey() string {
	return s.apiKey
}

// TextToSpeech 文字转语音 (TTS)
func (s *AIService) TextToSpeech(text string, characterName string) ([]byte, error) {
	start := time.Now()

	// 限制文本长度，提高TTS速度
	if len([]rune(text)) > 100 {
		text = string([]rune(text)[:100]) + "..."
	}

	// 根据角色选择音色和语速
	voiceType := s.getVoiceTypeForCharacter(characterName)
	speedRatio := s.getSpeedRatioForCharacter(characterName)

	req := TTSRequest{
		Audio: struct {
			VoiceType  string  `json:"voice_type"`
			Encoding   string  `json:"encoding"`
			SpeedRatio float64 `json:"speed_ratio"`
		}{
			VoiceType:  voiceType,
			Encoding:   "mp3",
			SpeedRatio: speedRatio,
		},
		Request: struct {
			Text      string `json:"text"`
			MaxTokens int    `json:"max_tokens,omitempty"`
		}{
			Text:      text,
			MaxTokens: 0, // 移除MaxTokens限制
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TTS request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.baseURL+"/voice/tts", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create TTS request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 为TTS创建专门的客户端，合理超时时间
	ttsClient := &http.Client{
		Timeout: 5 * time.Second, // 根据实际性能调整，TTS耗时约3.4秒
	}

	resp, err := ttsClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("TTS HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TTS API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var ttsResp TTSResponse
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read TTS response: %w", err)
	}

	if err := json.Unmarshal(respBody, &ttsResp); err != nil {
		return nil, fmt.Errorf("failed to parse TTS response: %w", err)
	}

	// 解码base64音频数据
	audioData, err := base64.StdEncoding.DecodeString(ttsResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode TTS audio data: %w", err)
	}

	fmt.Printf("TTS耗时: %v\n", time.Since(start))
	return audioData, nil
}

// getVoiceTypeForCharacter 根据角色名称获取对应的音色
func (s *AIService) getVoiceTypeForCharacter(characterName string) string {
	switch characterName {
	case "林黛玉":
		return "qiniu_zh_female_wwxkjx" // 温婉学科讲师 - 温婉细腻，符合林黛玉的气质
	case "孙悟空":
		return "qiniu_zh_male_mzjsxg" // 名著角色猴哥 - 专门为孙悟空设计的音色
	case "李白":
		return "qiniu_zh_male_gzjjxb" // 古装剧教学版 - 成熟正式，符合李白的豪迈洒脱
	case "赫敏·格兰杰", "赫敏":
		return "qiniu_zh_female_ljfdxx" // 知性教学女教师 - 清晰理性，有条理
	default:
		return "qiniu_zh_female_xyqxxj" // 校园清新学姐 - 默认女声
	}
}

// getSpeedRatioForCharacter 根据角色名称获取对应的语速
func (s *AIService) getSpeedRatioForCharacter(characterName string) float64 {
	switch characterName {
	case "林黛玉":
		return 0.9 // 语速稍慢，符合忧郁细腻的特点
	case "孙悟空":
		return 1.1 // 语速稍快，符合活泼有力的特点
	case "李白":
		return 1.0 // 语速正常，符合豪迈洒脱的特点
	case "赫敏·格兰杰", "赫敏":
		return 1.0 // 语速较快，符合清晰理性的特点
	default:
		return 1.0 // 默认正常语速
	}
}

// processEmojiMessages 处理表情消息，添加表情识别提示
func (s *AIService) processEmojiMessages(messages []Message) []Message {
	for i, msg := range messages {
		if msg.Role == "user" && s.isEmojiMessage(msg.Content) {
			// 为用户消息添加表情识别提示
			messages[i].Content = fmt.Sprintf("用户发送了表情：%s\n请根据这个表情和你的角色性格，自然地回应。", msg.Content)
		}
	}
	return messages
}

// isEmojiMessage 检查消息是否只包含表情
func (s *AIService) isEmojiMessage(content string) bool {
	// 简单的表情检测：如果消息只包含表情符号（Unicode表情范围）
	emojiRanges := []struct {
		start, end rune
	}{
		{0x1F600, 0x1F64F}, // Emoticons
		{0x1F300, 0x1F5FF}, // Misc Symbols and Pictographs
		{0x1F680, 0x1F6FF}, // Transport and Map
		{0x1F1E0, 0x1F1FF}, // Regional indicator symbols
		{0x2600, 0x26FF},   // Misc symbols
		{0x2700, 0x27BF},   // Dingbats
		{0xFE00, 0xFE0F},   // Variation Selectors
		{0x1F900, 0x1F9FF}, // Supplemental Symbols and Pictographs
		{0x1F018, 0x1F0F5}, // Playing cards
		{0x1F200, 0x1F2FF}, // Enclosed characters
	}

	for _, r := range content {
		isEmoji := false
		for _, emojiRange := range emojiRanges {
			if r >= emojiRange.start && r <= emojiRange.end {
				isEmoji = true
				break
			}
		}
		if !isEmoji && r != ' ' && r != '\n' && r != '\t' {
			return false
		}
	}
	return len(content) > 0
}

// IsEmojiMessage 检查消息是否只包含表情（公开方法）
func (s *AIService) IsEmojiMessage(content string) bool {
	return s.isEmojiMessage(content)
}

// convertPCMToWAV 将PCM数据转换为WAV格式
func (s *AIService) convertPCMToWAV(pcmData []byte) ([]byte, error) {
	// WAV文件头结构
	const (
		sampleRate    = 16000 // 采样率
		channels      = 1     // 单声道
		bitsPerSample = 16    // 16位
	)

	// 计算数据大小
	dataSize := len(pcmData)
	fileSize := 36 + dataSize

	// 创建WAV文件头
	header := make([]byte, 44)

	// RIFF头
	copy(header[0:4], "RIFF")
	header[4] = byte(fileSize)
	header[5] = byte(fileSize >> 8)
	header[6] = byte(fileSize >> 16)
	header[7] = byte(fileSize >> 24)
	copy(header[8:12], "WAVE")

	// fmt chunk
	copy(header[12:16], "fmt ")
	header[16] = 16                   // fmt chunk size
	copy(header[20:22], []byte{1, 0}) // PCM format
	header[22] = byte(channels)
	header[24] = byte(sampleRate & 0xFF)
	header[25] = byte((sampleRate >> 8) & 0xFF)
	header[26] = byte((sampleRate >> 16) & 0xFF)
	header[27] = byte((sampleRate >> 24) & 0xFF)

	byteRate := sampleRate * channels * bitsPerSample / 8
	header[28] = byte(byteRate & 0xFF)
	header[29] = byte((byteRate >> 8) & 0xFF)
	header[30] = byte((byteRate >> 16) & 0xFF)
	header[31] = byte((byteRate >> 24) & 0xFF)

	blockAlign := channels * bitsPerSample / 8
	header[32] = byte(blockAlign)
	header[33] = byte(blockAlign >> 8)
	header[34] = byte(bitsPerSample)
	header[35] = byte(bitsPerSample >> 8)

	// data chunk
	copy(header[36:40], "data")
	header[40] = byte(dataSize)
	header[41] = byte(dataSize >> 8)
	header[42] = byte(dataSize >> 16)
	header[43] = byte(dataSize >> 24)

	// 合并头部和PCM数据
	wavData := make([]byte, len(header)+len(pcmData))
	copy(wavData, header)
	copy(wavData[len(header):], pcmData)

	return wavData, nil
}
