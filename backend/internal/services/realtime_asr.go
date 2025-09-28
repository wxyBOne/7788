package services

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// RealtimeASRClient 实时ASR客户端
type RealtimeASRClient struct {
	wsURL       string
	token       string
	conn        *websocket.Conn
	seq         int32
	isConnected bool
	audioConfig AudioConfig
}

// AudioConfig 音频配置
type AudioConfig struct {
	SampleRate int    `json:"sample_rate"`
	Channels   int    `json:"channel"`
	Bits       int    `json:"bits"`
	Format     string `json:"format"`
	Codec      string `json:"codec"`
}

// RealtimeASRRequest 实时ASR配置请求
type RealtimeASRRequest struct {
	User    UserInfo    `json:"user"`
	Audio   AudioConfig `json:"audio"`
	Request RequestInfo `json:"request"`
}

// UserInfo 用户信息
type UserInfo struct {
	UID string `json:"uid"`
}

// RequestInfo 请求信息
type RequestInfo struct {
	ModelName  string `json:"model_name"`
	EnablePunc bool   `json:"enable_punc"`
}

// RealtimeASRResponse 实时ASR响应
type RealtimeASRResponse struct {
	PayloadMsg interface{} `json:"payload_msg"`
	Seq        int32       `json:"payload_sequence,omitempty"`
	IsLast     bool        `json:"is_last_package,omitempty"`
}

// NewRealtimeASRClient 创建实时ASR客户端
func NewRealtimeASRClient(token string) *RealtimeASRClient {
	return &RealtimeASRClient{
		wsURL: "wss://openai.qiniu.com/v1/voice/asr",
		token: token,
		seq:   1,
		audioConfig: AudioConfig{
			SampleRate: 16000,
			Channels:   1,
			Bits:       16,
			Format:     "pcm",
			Codec:      "raw",
		},
	}
}

// 协议常量
const (
	PROTOCOL_VERSION = 0x01
	HEADER_SIZE      = 0x01

	// Message Types
	FULL_CLIENT_REQUEST   = 0x01
	AUDIO_ONLY_REQUEST    = 0x02
	FULL_SERVER_RESPONSE  = 0x09
	SERVER_ACK            = 0x0B
	SERVER_ERROR_RESPONSE = 0x0F

	// Message Type Specific Flags
	NO_SEQUENCE  = 0x00
	POS_SEQUENCE = 0x01
	NEG_SEQUENCE = 0x02

	// Serialization and Compression
	JSON_SERIALIZATION = 0x01
	NO_COMPRESSION     = 0x00
	GZIP_COMPRESSION   = 0x01
)

// generateHeader 生成协议头
func (c *RealtimeASRClient) generateHeader(messageType, flags, serial, compress byte) []byte {
	header := make([]byte, 4)
	header[0] = (PROTOCOL_VERSION << 4) | HEADER_SIZE
	header[1] = (messageType << 4) | flags
	header[2] = (serial << 4) | compress
	header[3] = 0x00
	return header
}

// generateSequence 生成序列号
func (c *RealtimeASRClient) generateSequence(seq int32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(seq))
	return buf
}

// generatePayloadLength 生成负载长度
func (c *RealtimeASRClient) generatePayloadLength(length int) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(length))
	return buf
}

// Connect 连接到WebSocket
func (c *RealtimeASRClient) Connect() error {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+c.token)

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	fmt.Printf("正在连接到七牛云ASR: %s\n", c.wsURL)
	conn, _, err := dialer.Dial(c.wsURL, headers)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	fmt.Println("WebSocket连接已建立")
	c.conn = conn
	c.isConnected = true

	// 发送配置信息
	fmt.Println("发送ASR配置信息...")
	err = c.sendConfig()
	if err != nil {
		return fmt.Errorf("failed to send config: %w", err)
	}

	// 等待配置响应
	fmt.Println("等待配置响应...")
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, _, err = conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read config response: %w", err)
	}

	fmt.Println("配置响应已接收")
	// 重置读取超时
	conn.SetReadDeadline(time.Time{})

	return nil
}

// sendConfig 发送配置信息
func (c *RealtimeASRClient) sendConfig() error {
	req := RealtimeASRRequest{
		User: UserInfo{
			UID: "realtime-asr-client",
		},
		Audio: c.audioConfig,
		Request: RequestInfo{
			ModelName:  "asr",
			EnablePunc: true,
		},
	}

	// JSON序列化
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// GZIP压缩
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	_, err = gzWriter.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to compress config: %w", err)
	}
	gzWriter.Close()

	compressedData := buf.Bytes()

	// 构造消息
	header := c.generateHeader(FULL_CLIENT_REQUEST, POS_SEQUENCE, JSON_SERIALIZATION, GZIP_COMPRESSION)
	sequence := c.generateSequence(c.seq)
	payloadLength := c.generatePayloadLength(len(compressedData))

	message := append(header, sequence...)
	message = append(message, payloadLength...)
	message = append(message, compressedData...)

	// 发送消息
	err = c.conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		return fmt.Errorf("failed to send config: %w", err)
	}

	c.seq++
	return nil
}

// SendAudioChunk 发送音频分片
func (c *RealtimeASRClient) SendAudioChunk(audioData []byte) error {
	if !c.isConnected {
		return fmt.Errorf("WebSocket not connected")
	}

	fmt.Printf("发送音频数据: %d bytes\n", len(audioData))

	// 音频数据不需要GZIP压缩，直接发送PCM数据
	// 构造消息
	header := c.generateHeader(AUDIO_ONLY_REQUEST, POS_SEQUENCE, JSON_SERIALIZATION, NO_COMPRESSION)
	sequence := c.generateSequence(c.seq)
	payloadLength := c.generatePayloadLength(len(audioData))

	message := append(header, sequence...)
	message = append(message, payloadLength...)
	message = append(message, audioData...)

	// 发送消息
	err := c.conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		// 检查是否是连接错误
		if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) ||
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			c.isConnected = false
		}
		return fmt.Errorf("failed to send audio chunk: %w", err)
	}

	fmt.Printf("音频数据发送成功，序列号: %d\n", c.seq)
	c.seq++
	return nil
}

// ReadResponse 读取响应
func (c *RealtimeASRClient) ReadResponse() (string, error) {
	if !c.isConnected {
		return "", fmt.Errorf("WebSocket not connected")
	}

	_, data, err := c.conn.ReadMessage()
	if err != nil {
		// 检查是否是连接错误
		if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) ||
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			c.isConnected = false
		}
		return "", fmt.Errorf("failed to read message: %w", err)
	}

	return c.parseResponse(data)
}

// parseResponse 解析响应
func (c *RealtimeASRClient) parseResponse(data []byte) (string, error) {
	if len(data) < 4 {
		return "", fmt.Errorf("invalid response data")
	}

	// 解析头部
	headerSize := data[0] & 0x0F
	messageType := data[1] >> 4
	messageTypeSpecificFlags := data[1] & 0x0F
	serializationMethod := data[2] >> 4
	messageCompression := data[2] & 0x0F

	payload := data[headerSize*4:]

	// 处理序列号
	if messageTypeSpecificFlags&0x01 != 0 {
		if len(payload) < 4 {
			return "", fmt.Errorf("invalid payload")
		}
		payload = payload[4:] // 跳过序列号
	}

	// 处理负载大小
	if messageType == FULL_SERVER_RESPONSE && len(payload) >= 4 {
		payloadSize := binary.BigEndian.Uint32(payload[:4])
		payload = payload[4 : 4+payloadSize]
	}

	// 解压缩
	if messageCompression == GZIP_COMPRESSION {
		gzReader, err := gzip.NewReader(bytes.NewReader(payload))
		if err != nil {
			return "", fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()

		decompressed, err := io.ReadAll(gzReader)
		if err != nil {
			return "", fmt.Errorf("failed to decompress: %w", err)
		}
		payload = decompressed
	}

	// 解析JSON
	if serializationMethod == JSON_SERIALIZATION {
		var result map[string]interface{}
		err := json.Unmarshal(payload, &result)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		// 提取文本 - 支持多种响应格式
		// 格式1: result.text
		if resultMsg, ok := result["result"].(map[string]interface{}); ok {
			if text, ok := resultMsg["text"].(string); ok {
				return text, nil
			}
		}

		// 格式2: payload_msg.result.text (官方示例格式)
		if payloadMsg, ok := result["payload_msg"].(map[string]interface{}); ok {
			if resultMsg, ok := payloadMsg["result"].(map[string]interface{}); ok {
				if text, ok := resultMsg["text"].(string); ok {
					return text, nil
				}
			}
		}

		// 格式3: 直接是字符串
		if text, ok := result["payload_msg"].(string); ok {
			return text, nil
		}
	}

	return "", nil
}

// Close 关闭连接
func (c *RealtimeASRClient) Close() error {
	c.isConnected = false
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsConnected 检查连接状态
func (c *RealtimeASRClient) IsConnected() bool {
	return c.isConnected
}

// SetReadDeadline 设置读取超时
func (c *RealtimeASRClient) SetReadDeadline(t time.Time) error {
	if c.conn != nil {
		return c.conn.SetReadDeadline(t)
	}
	return fmt.Errorf("WebSocket not connected")
}
