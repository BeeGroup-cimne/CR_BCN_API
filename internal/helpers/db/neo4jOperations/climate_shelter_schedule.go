package neo4jOperations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func climateShelterFilterConditions(selectedCharacteristics, selectedBuildingType []string) string {
	var conds []string
	if len(selectedCharacteristics) != 0 {
		conds = append(conds, "ALL(attr IN $selectedCharacteristics WHERE n[attr] = \"True\")")
	}
	if len(selectedBuildingType) != 0 {
		conds = append(conds, "n.ns1__type IN $selectedBuildingType")
	}
	if len(conds) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(conds, " AND ") + " "
}

func selectedKeys(m map[string]bool) []string {
	keys := []string{}
	for k, v := range m {
		if v {
			keys = append(keys, k)
		}
	}
	return keys
}

func cleanScheduleJson(raw string) string {
	cleaned := strings.ReplaceAll(raw, "'", `"`)
	return strings.ReplaceAll(cleaned, "&apos;", "'")
}

func parseScheduleList(raw string) ([]interface{}, error) {
	if raw == "" {
		return nil, nil
	}
	dec := json.NewDecoder(bytes.NewReader([]byte(cleanScheduleJson(raw))))
	dec.UseNumber()

	var outer []interface{}
	if err := dec.Decode(&outer); err != nil {
		return nil, err
	}
	var flat []interface{}
	for _, o := range outer {
		if inner, ok := o.([]interface{}); ok {
			flat = append(flat, inner...)
		}
	}
	return flat, nil
}

type shelterPeriod struct {
	MonthStart int `json:"month_start"`
	DayStart   int `json:"day_start"`
	MonthEnd   int `json:"month_end"`
	DayEnd     int `json:"day_end"`
}

type shelterOpening struct {
	HourStart any `json:"hour_start"` // "09:00" o número
	HourEnd   any `json:"hour_end"`
}

type shelterProfile struct {
	Days         []any            `json:"days"`
	OpeningHours []shelterOpening `json:"opening_hours"`
}

type shelterSchedule struct {
	Periods       []shelterPeriod  `json:"periods"`
	WeeklyProfile []shelterProfile `json:"weekly_profile"`
}

func shelterIsOpenAt(raw string, t time.Time) (bool, error) {
	var schedules [][]shelterSchedule
	if err := json.Unmarshal([]byte(cleanScheduleJson(raw)), &schedules); err != nil {
		return false, err
	}

	weekday := t.Weekday().String()
	now := t.Hour()*60 + t.Minute()

	for _, outer := range schedules {
		for _, s := range outer {
			if !inAnyShelterPeriod(s.Periods, t) {
				continue
			}
			for _, prof := range s.WeeklyProfile {
				if !containsWeekday(prof.Days, weekday) {
					continue
				}
				for _, o := range prof.OpeningHours {
					start, ok1 := toMinutes(o.HourStart)
					end, ok2 := toMinutes(o.HourEnd)
					if ok1 && ok2 && inOpeningHours(start, end, now) {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

func inAnyShelterPeriod(periods []shelterPeriod, t time.Time) bool {
	if len(periods) == 0 {
		return true
	}
	md := int(t.Month())*100 + t.Day()
	for _, p := range periods {
		start := p.MonthStart*100 + p.DayStart
		end := p.MonthEnd*100 + p.DayEnd
		if start <= end {
			if md >= start && md <= end {
				return true
			}
		} else if md >= start || md <= end {
			return true
		}
	}
	return false
}

func inOpeningHours(start, end, now int) bool {
	if start <= end {
		return now >= start && now <= end
	}
	return now >= start || now <= end
}

func toMinutes(v any) (int, bool) {
	switch h := v.(type) {
	case float64:
		return int(h * 60), true
	case string:
		parts := strings.Split(strings.TrimSpace(h), ":")
		hh, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, false
		}
		mm := 0
		if len(parts) > 1 {
			if mm, err = strconv.Atoi(parts[1]); err != nil {
				return 0, false
			}
		}
		return hh*60 + mm, true
	}
	return 0, false
}

func containsWeekday(days []any, weekday string) bool {
	for _, d := range days {
		if strings.EqualFold(fmt.Sprint(d), weekday) {
			return true
		}
	}
	return false
}
