package ansutil

import (
	"fmt"
	app "go-rest-api-boilerplate/internal"
	"go-rest-api-boilerplate/pkg/imgUploader"
	"go-rest-api-boilerplate/types"
	"log/slog"
	"strings"
	"sync"
)

const (
	ImgHost = "https://i.ibb.co/"
)

type AnswersHandler interface {
	Handle(answers []types.AnswerItem, i imgUploader.IUploadImage) ([]types.AnswerItem, error)
}

type AnswerImg struct{}

func (a *AnswerImg) Handle(answers []types.AnswerItem, uploader imgUploader.IUploadImage) ([]types.AnswerItem, error) {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)  // Channel để nhận lỗi đầu tiên
	doneCh := make(chan struct{}) // Channel để gửi tín hiệu hủy bỏ
	var once sync.Once            // Đảm bảo hủy bỏ chỉ được thực hiện một lần

	for i := range answers {
		wg.Add(1)
		go func(item *types.AnswerItem) { // Truyền con trỏ để thay đổi giá trị gốc
			defer wg.Done()

			// Kiểm tra xem có tín hiệu hủy bỏ không
			select {
			case <-doneCh:
				return // Dừng goroutine nếu nhận được tín hiệu hủy bỏ
			default:
				if item.Ans != "" {
					// Kiểm tra xem item.Ans đã là URL hợp lệ chưa
					if strings.HasPrefix(item.Ans, ImgHost) {
						// Nếu đã là URL hợp lệ, bỏ qua bước upload
						return
					}

					// Nếu chưa là URL hợp lệ, thực hiện upload
					result, err := uploader.Upload(item.Ans)
					if err != nil {
						// Gửi lỗi vào channel và đóng channel hủy bỏ
						once.Do(func() {
							errCh <- fmt.Errorf("error uploading image %s: %v", item.Ans, err)
							close(doneCh) // Gửi tín hiệu hủy bỏ đến các goroutine khác
						})
						return
					}
					item.Ans = result.Url
				}
			}
		}(&answers[i]) // Truyền địa chỉ của phần tử trong slice
	}

	// Đợi tất cả các goroutine hoàn thành
	wg.Wait()

	// Kiểm tra xem có lỗi nào không
	select {
	case err := <-errCh:
		slog.Error("Upload error:", err)
		return nil, app.ErrInternalServerError
	default:
		// Nếu không có lỗi, tiếp tục xử lý
	}

	return answers, nil
}

type AnswerText struct{}

func (a *AnswerText) Handle(answers []types.AnswerItem, _ imgUploader.IUploadImage) ([]types.AnswerItem, error) {
	return answers, nil
}

// GetHandles returns interface
func GetHandler(t string) AnswersHandler {
	switch t {
	case "IMAGE":
		return &AnswerImg{}
	case "TEXT":
		return &AnswerText{}
	default:
		return nil
	}
}
