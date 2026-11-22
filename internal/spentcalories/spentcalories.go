package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	splittedData := strings.Split(data, ",")
	if len(splittedData) != 3 {
		return 0, "", 0, errors.New("некорретное число данных, ожидается формат: <число шагов>,<вид деятольности>,<продолжительность>")
	}
	steps, err := strconv.Atoi(splittedData[0])
	if err != nil {
		return 0, "", 0, err
	}
	trainingType := splittedData[1]
	duration, err := time.ParseDuration(splittedData[2])
	if err != nil {
		return 0, "", 0, err
	}
	return steps, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient //это все можно было бы засунуть в одно большое выражение,
	distance := stepLength * float64(steps)      //но мне показалось, что приятнее выглядит,
	distanceMInKm := distance / mInKm            //если все разбить по переменным ;)
	return distanceMInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	averageSpeed := distance / duration.Hours()
	return averageSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	distance := distance(steps, height)
	meanSpeed := meanSpeed(steps, weight, duration)
	var calories float64
	switch trainingType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	trainingInfo := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType,
		duration,
		distance,
		meanSpeed,
		calories,
	)
	return trainingInfo, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("аргументы должны быть больше 0")
	}
	meanSpeed := meanSpeed(steps, height, duration)
	calories := (duration.Minutes() * meanSpeed * weight) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	runningSpentCalories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	walkingSpentCalories := walkingCaloriesCoefficient * runningSpentCalories
	return walkingSpentCalories, nil
}
