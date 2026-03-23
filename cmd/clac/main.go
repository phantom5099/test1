package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorGray   = "\033[90m"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	history := make([]string, 0, 16)
	lastResult := 0.0
	hasLastResult := false

	printWelcome()

	for {
		fmt.Print(colorCyan + "calc> " + colorReset)
		line, err := reader.ReadString('\n')
		if err != nil {
			line = strings.TrimSpace(line)
			if line == "" {
				return
			}
		} else {
			line = strings.TrimSpace(line)
		}

		if line == "" {
			continue
		}

		switch strings.ToLower(line) {
		case "exit", "quit":
			fmt.Println(colorYellow + "已退出 clac" + colorReset)
			return
		case "help":
			printHelp()
			continue
		case "history":
			printHistory(history)
			continue
		case "clear":
			printClear()
			continue
		case "ans":
			if !hasLastResult {
				fmt.Println(colorGray + "暂无上一条结果" + colorReset)
				continue
			}
			fmt.Printf(colorGreen+"上次结果: %s\n"+colorReset, formatResult(lastResult))
			continue
		}

		if hasLastResult {
			line = replaceAns(line, lastResult)
		} else if containsAns(line) {
			fmt.Println(colorRed + "错误: 当前还没有上一条结果，无法使用 ans" + colorReset)
			continue
		}

		result, err := calculate(line)
		if err != nil {
			fmt.Printf(colorRed+"错误: %v\n"+colorReset, err)
			continue
		}

		hasLastResult = true
		lastResult = result
		history = append(history, fmt.Sprintf("%s = %s", line, formatResult(result)))
		fmt.Printf(colorGreen+"结果: %s\n"+colorReset, formatResult(result))
	}
}

func printWelcome() {
	fmt.Println(colorCyan + "╔══════════════════════════════════════╗")
	fmt.Println("║          Go clac 增强计算器         ║")
	fmt.Println("╚══════════════════════════════════════╝" + colorReset)
	fmt.Println("支持基础/扩展运算，输入 help 查看完整说明")
	fmt.Println("示例: 12 + 3, 10 xor 3, 27 root 3, sqrt 81")
	fmt.Println("命令: help / history / clear / ans / exit")
}

func printHelp() {
	fmt.Println(colorYellow + "使用说明:" + colorReset)
	fmt.Println("  一、双参数格式: 数字 运算符 数字")
	fmt.Println("     例如: 12 + 3")
	fmt.Println("     支持: +  -  *  /  %  pow  ^  xor  root  max  min")
	fmt.Println("     说明:")
	fmt.Println("       a xor b   -> 按位异或(仅整数)")
	fmt.Println("       a root b  -> a 的 b 次方根，例如 27 root 3 = 3")
	fmt.Println("       a ^ b     -> a 的 b 次方")
	fmt.Println("  二、单参数格式: 运算符 数字")
	fmt.Println("     支持: sqrt  cbrt  abs")
	fmt.Println("     例如: sqrt 81, cbrt 27, abs -9")
	fmt.Println("  三、命令:")
	fmt.Println("     help    查看帮助")
	fmt.Println("     history 查看历史记录")
	fmt.Println("     clear   清屏")
	fmt.Println("     ans     查看上一条结果")
	fmt.Println("     exit    退出程序")
	fmt.Println("  四、支持在表达式中使用 ans")
	fmt.Println("     例如: ans * 2, ans xor 7")
}

func printHistory(history []string) {
	if len(history) == 0 {
		fmt.Println(colorGray + "暂无历史记录" + colorReset)
		return
	}

	fmt.Println(colorYellow + "历史记录:" + colorReset)
	for i, item := range history {
		fmt.Printf("  %d. %s\n", i+1, item)
	}
}

func printClear() {
	for i := 0; i < 30; i++ {
		fmt.Println()
	}
	printWelcome()
}

func containsAns(input string) bool {
	return strings.Contains(strings.ToLower(input), "ans")
}

func replaceAns(input string, lastResult float64) string {
	value := formatResult(lastResult)
	replacer := strings.NewReplacer("ans", value, "ANS", value, "Ans", value)
	return replacer.Replace(input)
}

func formatResult(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

func calculate(input string) (float64, error) {
	fields := strings.Fields(input)

	if len(fields) == 2 {
		return calculateUnary(fields[0], fields[1])
	}

	if len(fields) != 3 {
		return 0, fmt.Errorf("请输入: 数字 运算符 数字，或 运算符 数字，例如 12 + 3 / sqrt 81")
	}

	left, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("左操作数无效: %s", fields[0])
	}

	right, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return 0, fmt.Errorf("右操作数无效: %s", fields[2])
	}

	switch strings.ToLower(fields[1]) {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("除数不能为 0")
		}
		return left / right, nil
	case "%":
		if right == 0 {
			return 0, fmt.Errorf("取模时除数不能为 0")
		}
		return math.Mod(left, right), nil
	case "^", "pow":
		return math.Pow(left, right), nil
	case "root":
		return nthRoot(left, right)
	case "xor":
		li, err := parseIntegerOperand(fields[0])
		if err != nil {
			return 0, fmt.Errorf("异或运算左操作数必须是整数")
		}
		ri, err := parseIntegerOperand(fields[2])
		if err != nil {
			return 0, fmt.Errorf("异或运算右操作数必须是整数")
		}
		return float64(li ^ ri), nil
	case "max":
		return math.Max(left, right), nil
	case "min":
		return math.Min(left, right), nil
	default:
		return 0, fmt.Errorf("不支持的运算符: %s", fields[1])
	}
}

func calculateUnary(op string, valueText string) (float64, error) {
	value, err := strconv.ParseFloat(valueText, 64)
	if err != nil {
		return 0, fmt.Errorf("操作数无效: %s", valueText)
	}

	switch strings.ToLower(op) {
	case "sqrt":
		if value < 0 {
			return 0, fmt.Errorf("负数不能开平方根")
		}
		return math.Sqrt(value), nil
	case "cbrt":
		return math.Cbrt(value), nil
	case "abs":
		return math.Abs(value), nil
	default:
		return 0, fmt.Errorf("不支持的一元运算: %s", op)
	}
}

func nthRoot(value float64, degree float64) (float64, error) {
	if degree == 0 {
		return 0, fmt.Errorf("根指数不能为 0")
	}
	if value < 0 {
		degreeInt, ok := exactInteger(degree)
		if !ok || degreeInt%2 == 0 {
			return 0, fmt.Errorf("负数只支持奇数次方根")
		}
		return -math.Pow(-value, 1/degree), nil
	}
	return math.Pow(value, 1/degree), nil
}

func parseIntegerOperand(text string) (int64, error) {
	v, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, err
	}
	i, ok := exactInteger(v)
	if !ok {
		return 0, fmt.Errorf("不是整数")
	}
	return i, nil
}

func exactInteger(v float64) (int64, bool) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	i := int64(v)
	return i, float64(i) == v
}
