package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gavink97/pgn-tools/internal/global"
	"github.com/gavink97/pgn-tools/internal/types"
)

type QueryCondition struct {
	Key   string
	Op    string
	Value string
}

type Query struct {
	Conditions []QueryCondition
}

func ParseQuery(keys string) (*Query, error) {
	keyValues := strings.Split(keys, ",")
	query := &Query{}
	operators := []string{">=", "<=", "!=", "=", ">", "<"}

	var op, key, value string

	for _, kv := range keyValues {
		kv = strings.TrimSpace(kv)
		if kv == "" {
			continue
		}

		for _, operator := range operators {
			if strings.Contains(kv, operator) {
				parts := strings.SplitN(kv, operator, 2)
				if len(parts) == 2 {
					key = strings.TrimSpace(parts[0])
					op = operator
					value = strings.TrimSpace(parts[1])
					break
				}
			}
		}

		key = strings.ToLower(key)

		if key == "" || op == "" {
			return nil, fmt.Errorf("invalid key value pair: %s", kv)
		}

		query.Conditions = append(query.Conditions, QueryCondition{
			Key:   key,
			Op:    op,
			Value: value,
		})

		global.Logger.Debug(fmt.Sprintf("Parsed condition: %s %s %s", key, op, value))
	}

	return query, nil
}

func (q *Query) Match(game *types.Game) (bool, error) {
	for _, condition := range q.Conditions {
		matches, err := condition.Evaluate(game)
		if err != nil {
			return false, err
		}

		if !matches {
			return false, nil
		}
	}

	return true, nil
}

func (c *QueryCondition) Evaluate(game *types.Game) (bool, error) {
	computedFunc, exists := computedFields[strings.ToLower(c.Key)]
	if exists {
		value, err := computedFunc(game, c)
		if err != nil {
			return false, fmt.Errorf("an unexpected error occured: %v", err)
		}

		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			return c.EvaluateString(v)
		case int:
			return c.EvaluateInt(int64(v))
		case int64:
			return c.EvaluateInt(v)
		default:
			return false, fmt.Errorf("unsupported computed field type: %v", value)
		}
	}

	gameValue := reflect.ValueOf(game).Elem()
	field, found := FindField(gameValue, c.Key)

	if !found {
		return false, fmt.Errorf("unknown field: %s", c.Key)
	}

	switch field.Kind() {
	case reflect.String:
		return c.EvaluateString(field.String())
	case reflect.Int:
		return c.EvaluateInt(field.Int())
	default:
		return false, fmt.Errorf("unsupport field type: %s", field.Kind())
	}
}

func (c *QueryCondition) EvaluateString(value string) (bool, error) {
	value = strings.ToLower(value)
	cValue := strings.ToLower(c.Value)

	switch c.Op {
	case "=":
		return strings.Contains(value, cValue), nil
	case "!=":
		return !strings.Contains(value, cValue), nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", c.Op)
	}
}

func (c *QueryCondition) EvaluateInt(value int64) (bool, error) {
	intValue, err := strconv.ParseInt(c.Value, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid integer value: %v", c.Value)
	}

	switch c.Op {
	case "=":
		return value == intValue, nil
	case "!=":
		return value != intValue, nil
	case ">":
		return value > intValue, nil
	case "<":
		return value < intValue, nil
	case ">=":
		return value >= intValue, nil
	case "<=":
		return value <= intValue, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", c.Op)
	}
}

func FindField(value reflect.Value, fieldName string) (reflect.Value, bool) {
	fieldName = strings.ToLower(fieldName)
	vType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := vType.Field(i)

		if strings.ToLower(structField.Name) == fieldName {
			return field, true
		}
	}
	return reflect.Value{}, false
}

var computedFields = map[string]func(*types.Game, *QueryCondition) (any, error){
	"elo": func(g *types.Game, qc *QueryCondition) (any, error) {
		if g.WhiteElo > 0 && g.BlackElo > 0 {
			return min(g.WhiteElo, g.BlackElo), nil
		} else {
			return 0, nil
		}
	},
	"player": func(g *types.Game, qc *QueryCondition) (any, error) {
		if qc.Op == "!=" {
			for _, player := range []string{g.White, g.Black} {
				player = strings.ToLower(player)
				value := strings.ToLower(qc.Value)

				if strings.Contains(player, value) {
					return false, nil
				}
			}

			return true, nil
		}

		white, err := qc.EvaluateString(g.White)
		if err != nil {
			return false, err
		}

		black, err := qc.EvaluateString(g.Black)
		if err != nil {
			return false, err
		}

		return white || black, nil
	},
}
