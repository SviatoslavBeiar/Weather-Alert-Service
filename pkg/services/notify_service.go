package services

import (
	"fmt"
	models2 "myapp/pkg/models"
	"myapp/pkg/utils"
	"regexp"
	"strconv"
	"strings"
)

var tempCondRe = regexp.MustCompile(
	`^temp\s*(<=|>=|<|>|==|=|!=)\s*([0-9]+(?:\.[0-9]+)?)$`,
)

func EvaluateAndNotify(sub models2.Subscription, weather models2.Weather) (bool, error) {
	cond := strings.TrimSpace(sub.Condition)
	var shouldSend bool

	if m := tempCondRe.FindStringSubmatch(cond); len(m) == 3 {
		op, thr := m[1], m[2]
		threshold, err := strconv.ParseFloat(thr, 64)
		if err != nil {
			return false, fmt.Errorf("invalid threshold %q: %v", thr, err)
		}
		switch op {
		case "<":
			shouldSend = weather.Temperature < threshold
		case "<=":
			shouldSend = weather.Temperature <= threshold
		case ">":
			shouldSend = weather.Temperature > threshold
		case ">=":
			shouldSend = weather.Temperature >= threshold
		case "=", "==":
			shouldSend = weather.Temperature == threshold
		case "!=":
			shouldSend = weather.Temperature != threshold
		default:
			// Теоретично сюди не зайде — regexp вже обмежує список
			return false, fmt.Errorf("unsupported operator %q", op)
		}
	} else if strings.EqualFold(cond, "rain") {
		shouldSend = strings.EqualFold(weather.Condition, "Rain")
	} else {
		return false, fmt.Errorf("unknown condition %q", cond)
	}

	if !shouldSend {
		return false, nil
	}

	subject := fmt.Sprintf("Weather Alert for %s", sub.City)
	body := fmt.Sprintf("Condition %s met: current temp %.1f°C", cond, weather.Temperature)
	if err := utils.SendEmail(sub.Email, subject, body); err != nil {
		return false, err
	}
	return true, nil
}
