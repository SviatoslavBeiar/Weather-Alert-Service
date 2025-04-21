package services_test

import (
	"errors"
	"testing"

	"myapp/pkg/models"
	"myapp/pkg/services"
	"myapp/pkg/utils"
)

func TestEvaluateAndNotify(t *testing.T) {
	orig := utils.SendEmail
	defer func() { utils.SendEmail = orig }()

	tests := []struct {
		name        string
		condition   string
		temp        float64
		weatherCond string
		emailErr    error
		wantSent    bool
		wantErr     bool
	}{
		{"LessThanTrue", "temp < 10", 5, "Sunny", nil, true, false},
		{"LessThanFalse", "temp < 0", 5, "Cloudy", nil, false, false},
		{"GreaterEqTrue", "temp >= 5", 5, "Sunny", nil, true, false},
		{"EqTrue", "temp == 5", 5, "Sunny", nil, true, false},
		{"NotEqTrue", "temp != 5", 6, "Sunny", nil, true, false},
		{"RainTrue", "rain", 20, "Rain", nil, true, false},
		{"RainFalse", "rain", 20, "Clear", nil, false, false},
		{"InvalidThreshold", "temp < abc", 5, "", nil, false, true},
		{"UnknownCondition", "snow", 0, "", nil, false, true},
		{"EmailError", "temp > 0", 10, "Sunny", errors.New("fail"), false, true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			called := false
			utils.SendEmail = func(_, _, _ string) error {
				called = true
				return tc.emailErr
			}

			sub := models.Subscription{Condition: tc.condition, Email: "a@b", City: "C"}
			w := models.Weather{Temperature: tc.temp, Condition: tc.weatherCond}

			sent, err := services.EvaluateAndNotify(sub, w)
			if (err != nil) != tc.wantErr {
				t.Fatalf("want err=%v, got %v", tc.wantErr, err)
			}
			if sent != tc.wantSent {
				t.Errorf("want sent=%v, got %v", tc.wantSent, sent)
			}
			if tc.wantSent && !called {
				t.Error("expected SendEmail to be called")
			}
		})
	}
}
