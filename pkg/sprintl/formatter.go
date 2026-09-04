package sprintl

import (
	"fmt"
)

// Package Sprintl
// TODO: Use reflect package in R and J functions so RF and
//       JF can be removed.
// TODO: Refactor, simplify, and tidy

type lineFormatter struct {
	typ            string
	delim          string
	args           []any
	values         []any
	valueArgs      [][]any
	generator      LineGenerator
	generatorCount int
}

func (lf lineFormatter) apply(line string) []string {
	switch lf.typ {
	case "F":
		return lf.applyF(line)
	case "R":
		return lf.applyR(line)
	case "RF":
		return lf.applyRF(line)
	case "J":
		return lf.applyJ(line)
	case "JF":
		return lf.applyJF(line)
	case "G":
		return lf.applyG(line)
	case "GN":
		return lf.applyGN(line)
	default:
		msg := fmt.Sprintf("Unknown format type '%s'", lf.typ)
		panic(msg)
	}
}

func (lf lineFormatter) applyF(
	templateLine string,
) []string {
	return []string{
		fmt.Sprintf(templateLine, lf.args...),
	}
}

func (lf lineFormatter) applyR(
	templateLine string,
) []string {
	lineCount := len(lf.values)
	lines := make([]string, lineCount, lineCount)

	for i, v := range lf.values {
		lines[i] = fmt.Sprintf(templateLine, v)
	}

	return lines
}

func (lf lineFormatter) applyRF(
	templateLine string,
) []string {
	lineCount := len(lf.valueArgs)
	lines := make([]string, lineCount, lineCount)

	for i, args := range lf.valueArgs {
		lines[i] = fmt.Sprintf(templateLine, args...)
	}

	return lines
}

func (lf lineFormatter) applyJ(
	templateLine string,
) []string {
	lineCount := len(lf.values)
	lines := make([]string, lineCount, lineCount)

	for i, v := range lf.values {
		if i > 0 {
			lines[i-1] += lf.delim
		}
		lines[i] = fmt.Sprintf(templateLine, v)
	}

	return lines
}

func (lf lineFormatter) applyJF(
	templateLine string,
) []string {
	lineCount := len(lf.valueArgs)
	lines := make([]string, lineCount, lineCount)

	for i, args := range lf.valueArgs {
		if i > 0 {
			lines[i-1] += lf.delim
		}
		lines[i] = fmt.Sprintf(templateLine, args...)
	}

	return lines
}

func (lf lineFormatter) applyG(
	templateLine string,
) []string {
	lines := []string{}

	for i := 0; true; i++ {
		line, ok := lf.generator(i, templateLine)
		if !ok {
			break
		}
		lines = append(lines, line)
	}

	return lines
}

func (lf lineFormatter) applyGN(
	templateLine string,
) []string {
	lines := []string{}

	for i := 0; i < lf.generatorCount; i++ {
		line, ok := lf.generator(i, templateLine)
		if ok {
			lines = append(lines, line)
		}
	}

	return lines
}
