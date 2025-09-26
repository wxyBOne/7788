package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"time"
)

type CaptchaService struct {
	width  int
	height int
}

type CaptchaData struct {
	Code     string `json:"code"`
	ImageURL string `json:"image_url"`
}

func NewCaptchaService() *CaptchaService {
	return &CaptchaService{
		width:  120,
		height: 40,
	}
}

// GenerateCaptcha 生成图片验证码
func (s *CaptchaService) GenerateCaptcha() (*CaptchaData, error) {
	// 生成4位随机字母
	code := s.generateRandomCode(4)

	// 创建图片
	img := image.NewRGBA(image.Rect(0, 0, s.width, s.height))

	// 填充背景色
	bgColor := color.RGBA{240, 240, 240, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 绘制验证码文字
	s.drawText(img, code)

	// 添加干扰线
	s.drawLines(img)

	// 添加噪点
	s.drawNoise(img)

	// 转换为base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	imageURL := fmt.Sprintf("data:image/png;base64,%s", imageBase64)

	return &CaptchaData{
		Code:     code,
		ImageURL: imageURL,
	}, nil
}

// generateRandomCode 生成随机字母代码
func (s *CaptchaService) generateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	code := make([]byte, length)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

// drawText 绘制文字
func (s *CaptchaService) drawText(img *image.RGBA, text string) {
	// 简单的文字绘制（实际项目中应该使用字体库）
	charWidth := s.width / len(text)
	charHeight := s.height - 10

	for i, char := range text {
		x := i*charWidth + 5
		y := charHeight/2 + 5

		// 随机颜色
		textColor := color.RGBA{
			uint8(rand.Intn(100) + 50),
			uint8(rand.Intn(100) + 50),
			uint8(rand.Intn(100) + 50),
			255,
		}

		// 绘制字符（简化版本，实际应该使用字体）
		s.drawChar(img, char, x, y, textColor)
	}
}

// drawChar 绘制单个字符（简化版本）
func (s *CaptchaService) drawChar(img *image.RGBA, char rune, x, y int, color color.RGBA) {
	// 这里简化处理，实际项目中应该使用字体库
	// 现在只是绘制一个简单的矩形代表字符
	for i := 0; i < 15; i++ {
		for j := 0; j < 20; j++ {
			if x+i < s.width && y+j < s.height {
				img.Set(x+i, y+j, color)
			}
		}
	}
}

// drawLines 绘制干扰线
func (s *CaptchaService) drawLines(img *image.RGBA) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		x1 := rand.Intn(s.width)
		y1 := rand.Intn(s.height)
		x2 := rand.Intn(s.width)
		y2 := rand.Intn(s.height)

		lineColor := color.RGBA{
			uint8(rand.Intn(200)),
			uint8(rand.Intn(200)),
			uint8(rand.Intn(200)),
			255,
		}

		s.drawLine(img, x1, y1, x2, y2, lineColor)
	}
}

// drawLine 绘制直线
func (s *CaptchaService) drawLine(img *image.RGBA, x1, y1, x2, y2 int, color color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1

	steps := max(abs(dx), abs(dy))
	if steps == 0 {
		return
	}

	xInc := float64(dx) / float64(steps)
	yInc := float64(dy) / float64(steps)

	x := float64(x1)
	y := float64(y1)

	for i := 0; i <= steps; i++ {
		if int(x) >= 0 && int(x) < s.width && int(y) >= 0 && int(y) < s.height {
			img.Set(int(x), int(y), color)
		}
		x += xInc
		y += yInc
	}
}

// drawNoise 绘制噪点
func (s *CaptchaService) drawNoise(img *image.RGBA) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		x := rand.Intn(s.width)
		y := rand.Intn(s.height)

		noiseColor := color.RGBA{
			uint8(rand.Intn(255)),
			uint8(rand.Intn(255)),
			uint8(rand.Intn(255)),
			255,
		}

		img.Set(x, y, noiseColor)
	}
}

// 辅助函数
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
