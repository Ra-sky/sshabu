package compare

import (
	"os"
	"bufio"
	"fmt"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
	Added ChangeType = iota
    Deleted
    Modified
    Moved
)

type ChangeType int

type Difference struct {
    lineNumber int
    line       string
    Added      bool
}

type Bites struct {
	length       int
	Content      []string
}

// External functions
func (bites *Bites) TakeBites(path string) {
    var lineArray []string

    file, err := os.Open(path)
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        lineArray = append(lineArray, line)
    }

    bites.Content = lineArray
    bites.length = len(lineArray)
}

// Internal functions
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func TransformDifferencesToReadableFormat(differences []Difference, firstBites Bites, secondBites Bites) []string {
    var result []string
    for index, line := range secondBites.Content {
        color := White
        resultStr := ""
        resultStr = fmt.Sprintf("%d: %s%s%s", index+1, color, line, White)
        for _, diff := range differences {
            if diff.lineNumber == index+1 {
                if diff.Added {
                    color = Green
                    resultStr = fmt.Sprintf("%d: %s%s%s", index+1, color, line, White)
                } else {
                    color = Red
                    resultStr = fmt.Sprintf("%d: %s%s\n    %s%s%s", index+1, color, diff.line ,Green, line, White)
                }
                break
            }
        }
        result = append(result, resultStr)
    }
    
    if len(result) < firstBites.length{
        for _, diff := range differences {
            if diff.lineNumber > len(result){
                color := Red
                resultStr := fmt.Sprintf("%d: %s%s%s", diff.lineNumber, color, diff.line, White)
                result = append(result, resultStr)
            }
        }
    }

    return result
}

func DiffBites(bites1, bites2 Bites) []Difference{
    var differences []Difference
    maxLen := 0
    for _, line := range bites1.Content {
        if len(line) > maxLen {
            maxLen = len(line)
        }
    }

    lcsMatrix := make([][]int, len(bites1.Content)+1)
    for i := range lcsMatrix {
        lcsMatrix[i] = make([]int, len(bites2.Content)+1)
    }

    // Построение матрицы LCS
    for i := 1; i <= len(bites1.Content); i++ {
        for j := 1; j <= len(bites2.Content); j++ {
            if bites1.Content[i-1] == bites2.Content[j-1] {
                lcsMatrix[i][j] = lcsMatrix[i-1][j-1] + 1
            } else {
                lcsMatrix[i][j] = max(lcsMatrix[i-1][j], lcsMatrix[i][j-1])
            }
        }
    }

    i, j := len(bites1.Content), len(bites2.Content)
    for i > 0 || j > 0 {
        if i > 0 && j > 0 && bites1.Content[i-1] == bites2.Content[j-1] {
            i--
            j--
        } else if j > 0 && (i == 0 || lcsMatrix[i][j-1] >= lcsMatrix[i-1][j]) {
            differences = append([]Difference{{lineNumber: j, line: bites2.Content[j-1], Added: true}}, differences...)
            j--
        } else if i > 0 && (j == 0 || lcsMatrix[i][j-1] < lcsMatrix[i-1][j]) {
            differences = append([]Difference{{lineNumber: i, line: bites1.Content[i-1], Added: false}}, differences...)
            i--
        }
    }

    return differences
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}





