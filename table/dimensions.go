package table

func (m *Model) recalculateWidth() {
	if m.targetTotalWidth != 0 {
		m.totalWidth = m.targetTotalWidth
	} else {
		total := 0

		for _, column := range m.columns {
			total += column.width
		}

		m.totalWidth = total + len(m.columns) + 1
	}

	updateColumnWidths(m.columns, m.targetTotalWidth)

	m.recalculateLastHorizontalColumn()
}

// Updates column width in-place.  This could be optimized but should be called
// very rarely so we prioritize simplicity over performance here.
func updateColumnWidths(cols []Column, totalWidth int) {
	totalFlexWidth := totalWidth - len(cols) - 1
	totalFlexFactor := 0
	flexGCD := 0

	for index, col := range cols {
		if !col.isFlex() {
			totalFlexWidth -= col.width
			cols[index].style = col.style.Width(col.width)
		} else {
			totalFlexFactor += col.flexFactor
			flexGCD = gcd(flexGCD, col.flexFactor)
		}
	}

	if totalFlexFactor == 0 {
		return
	}

	// We use the GCD here because otherwise very large values won't divide
	// nicely as ints
	totalFlexFactor /= flexGCD

	flexUnit := totalFlexWidth / totalFlexFactor
	leftoverWidth := totalFlexWidth % totalFlexFactor

	for index := range cols {
		if !cols[index].isFlex() {
			continue
		}

		width := flexUnit * (cols[index].flexFactor / flexGCD)

		if leftoverWidth > 0 {
			width++
			leftoverWidth--
		}

		if index == len(cols)-1 {
			width += leftoverWidth
			leftoverWidth = 0
		}

		width = max(width, 1)

		cols[index].width = width

		// Take borders into account for the actual style
		cols[index].style = cols[index].style.Width(width)
	}
}

func (m *Model) recalculateHeight() {
	var verticalChromeHeight int

	if m.headerVisible {
		verticalChromeHeight += 3
	}

	if m.hasFooter() {
		verticalChromeHeight += 2
	}

	// this accounts for the border at the bottom
	// when there's no footer
	verticalChromeHeight++

	m.maxVisibleRows = m.maxTotalHeight - verticalChromeHeight
	m.verticalScrollWindowStart = m.rowCursorIndex

	var newEnd = m.rowCursorIndex + m.maxVisibleRows - 1
	var lastRowIndex = m.TotalRows()

	if newEnd > lastRowIndex {
		m.verticalScrollWindowEnd = lastRowIndex - m.verticalScrollWindowStart - 1
	} else {
		m.verticalScrollWindowEnd = m.rowCursorIndex + min(m.maxVisibleRows, lastRowIndex) - 1
	}

	m.verticalScrollWindowYOffset = 0
}

func (m *Model) SetMaxTotalHeight(maxTotalHeight int) {
	// turn off paginated mode
	m.pageSize = 0
	m.rowCursorIndex = 0
	m.maxTotalHeight = maxTotalHeight
	m.recalculateHeight()
}

func (m *Model) SetMaxTotalWidth(maxTotalWidth int) {
	m.maxTotalWidth = maxTotalWidth

	m.recalculateWidth()
}
