package handler

import (
	"math"
)

const (
	minScore = 0
	maxScore = 100
	precision = 2
)

func Grade(subject1, subject2, subject3 int) (string, float64) {
	var grade string

	subject1GPA := calculateGPA(subject1)
	subject2GPA := calculateGPA(subject2)
	subject3GPA := calculateGPA(subject3)

	total := (subject1GPA + subject2GPA + subject3GPA) / 3
	total = ToFixed(total, precision)

	if subject1GPA == 0 || subject2GPA == 0 || subject3GPA == 0 {
		grade = "F"
		total = 0
	}

	if total <= 1.99 && total >= 1 {
		grade = "D"
	} else if total <= 2.99 && total >= 2 {
		grade = "C"
	} else if total <= 3.49 && total >= 3 {
		grade = "B"
	} else if total <= 3.99 && total >= 3.5 {
		grade = "A-"
	} else if total <= 4.99 && total >= 4 {
		grade = "A"
	} else if total == 5 {
		grade = "A+"
	} else if total <= 0.99 && total >= 0 {
		grade = "F"
	} else {
		grade = "Invalid"
	}

	return grade, total
}

func calculateGPA(score int) float64 {
	var GPA float64

	if score < minScore || score > maxScore {
		return 0
	}

	if score >= 80 {
		GPA = 5
	} else if score >= 70 {
		GPA = 4
	} else if score >= 60 {
		GPA = 3.5
	} else if score >= 50 {
		GPA = 3
	} else if score >= 40 {
		GPA = 2
	} else if score >= 33 {
		GPA = 1
	} else {
		GPA = 0
	}

	return GPA
}

func ToFixed(n float64, precision int) float64 {
	scale := math.Pow(10, float64(precision))
	return math.Round(n*scale) / scale
}
