package core

// FrameType defines the type of a VANTUN frame.
type FrameType uint8

const (
	// FrameTypeData is a data frame (0).
	FrameTypeData FrameType = 0
	// FrameTypePadding is a padding frame (1).
	FrameTypePadding FrameType = 1
	// FrameTypeTelemetry is a telemetry frame (2).
	FrameTypeTelemetry FrameType = 2
)