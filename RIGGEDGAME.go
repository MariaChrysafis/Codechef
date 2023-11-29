package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
)

func readInt(in *bufio.Reader) int {
	nStr, _ := in.ReadString('\n')
	nStr = strings.ReplaceAll(nStr, "\r", "")
	nStr = strings.ReplaceAll(nStr, "\n", "")
	n, _ := strconv.Atoi(nStr)
	return n
}

func readLineNumbs(in *bufio.Reader) []string {
	line, _ := in.ReadString('\n')
	line = strings.ReplaceAll(line, "\r", "")
	line = strings.ReplaceAll(line, "\n", "")
	numbs := strings.Split(line, " ")
	return numbs
}

func readArrInt(in *bufio.Reader) []int {
	numbs := readLineNumbs(in)
	arr := make([]int, len(numbs))
	for i, n := range numbs {
		val, _ := strconv.Atoi(n)
		arr[i] = val
	}
	return arr
}

func readArrInt64(in *bufio.Reader) []int64 {
	numbs := readLineNumbs(in)
	arr := make([]int64, len(numbs))
	for i, n := range numbs {
		val, _ := strconv.ParseInt(n, 10, 64)
		arr[i] = val
	}
	return arr
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	in := bufio.NewReader(os.Stdin)
	mod := int(1e9 + 7)
	mult := func(x int, y int) int {
		return int((int64(x) * int64(y)) % int64(mod))
	}
	add := func(x int, y int) int {
		return int((int64(x) + int64(y)) % int64(mod))
	}
	binPow := func(x int, y int) int {
		res := x
		ans := 1
		for y != 0 {
			if y%2 == 1 {
				ans = mult(ans, res)
			}
			res = mult(res, res)
			y /= 2
		}
		return ans
	}
	inv := func(x int) int {
		return binPow(x, mod-2)
	}
	tc := readInt(in)
	for t := 0; t < tc; t++ {
		n := readInt(in)
		arr := make([]int, n)
		adj := make([][]int, n)
		arr = readArrInt(in)
		for i := 0; i < n-1; i++ {
			a := readArrInt(in)
			x := a[0]
			y := a[1]
			x -= 1
			y -= 1
			adj[x] = append(adj[x], y)
			adj[y] = append(adj[y], x)
		}
		sub := make([]int, n) //number of guys in subtree
		var init func(int, int)
		init = func(curNode int, prevNode int) {
			sub[curNode] = 1
			for _, i := range adj[curNode] {
				if i != prevNode {
					init(i, curNode)
					sub[curNode] += sub[i]
				}
			}
		}
		init(0, 0)
		var brute func(int, int, int) int64
		brute = func(curNode int, prevNode int, d int) int64 { //return sum of depths
			ans := int64(d)
			for _, i := range adj[curNode] {
				if i != prevNode {
					ans += brute(i, curNode, d+1)
				}
			}
			return ans
		}
		weight := make([]int64, n)
		weight[0] = brute(0, 0, 0)
		var dfs func(int, int)
		dfs = func(curNode int, prevNode int) { //return sum of depths
			if curNode != prevNode {
				weight[curNode] = weight[prevNode]
				weight[curNode] -= int64(sub[curNode])
				weight[curNode] += int64(n - sub[curNode])
			}
			for _, i := range adj[curNode] {
				if i != prevNode {
					dfs(i, curNode)
				}
			}
		}
		dfs(0, 0)
		for i := 0; i < n; i++ {
			weight[i] = int64(n)*int64(n) - weight[i]
		}
		sort.Slice(arr, func(i int, j int) bool {
			return arr[i] < arr[j]
		})
		sort.Slice(weight, func(i int, j int) bool {
			return weight[i] > weight[j]
		})
		ans := 0
		for i := 0; i < len(weight); i++ {
			ans = add(ans, mult(int(weight[i]%int64(mod)), arr[i]))
		}
		fmt.Println(mult(ans, mult(inv(n), inv(n))))
	}
}
