package analysis

// represents a [start, end) range
type rangeInfo struct {
	start int
	end   int
}

func (r rangeInfo) isBefore(other rangeInfo) bool {
	return r.end <= other.start
}

func (r rangeInfo) isAfter(other rangeInfo) bool {
	return r.start >= other.end
}

func (r rangeInfo) isStartBefore(other rangeInfo) bool {
	return r.start < other.start
}

func (r rangeInfo) isEndBeforeEnd(other rangeInfo) bool {
	return r.end < other.end
}

func (r rangeInfo) startToStartDistance(other rangeInfo) int {
	return abs(other.start - r.start)
}

func (r rangeInfo) endToStartDistance(other rangeInfo) int {
	return abs(r.end - other.start)
}

func (r rangeInfo) endToEndDistance(other rangeInfo) int {
	return abs(r.end - other.end)
}

func (r rangeInfo) length() int {
	return r.end - r.start
}
