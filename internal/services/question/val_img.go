package question

type Image string

func (i Image) String() string {
	return string(i)

}

func NewImage(i string) (Image, error) {
	return Image(i), nil
}
