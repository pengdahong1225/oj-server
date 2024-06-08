package captcha

import "testing"

func TestGenerateImageCaptcha(t *testing.T) {
	// 调用
	id, b64s, err := GenerateImageCaptcha()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(id)
		t.Log(b64s)
	}
}

func TestVerifyImageCaptcha(t *testing.T) {
	id := "bn228CrYgRFxH2uU7UmJ"
	value := "cu9x"
	if VerifyImageCaptcha(id, value) {
		t.Log("true")
	} else {
		t.Log("false")
	}
}
