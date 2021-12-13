package models

func (r AuthRequest) IsValid() bool {
	return r.Amount.MajorUnits > 0
}

func (r CaptureRequest) IsValid() bool {
	return r.Amount.MajorUnits > 0
}

func (r RefundRequest) IsValid() bool {
	return r.Amount.MajorUnits > 0
}
