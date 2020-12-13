package main

import (
	"fmt"

	"github.com/kckecheng/pushgateway_pusher/parser"
)

func main() {
	// vdbenchfs := `07:32:31.019            1 1467.8  2.592  17.9 13.5   0.0    0.0  0.000` +
	// 	` 1467.8  2.592  0.00 183.4 183.47  131072   0.0  0.000   0.0  0.000   0.0  0.000  92.0` +
	// 	`  0.780  91.6  0.008   0.0  0.000`

	// p, err := parser.New("fields.yaml")
	// if err != nil {
	// 	panic(err)
	// }

	// for _, field := range p.Fields {
	// 	fmt.Println(field)
	// }

	// values, err := p.Parse(vdbenchfs)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, v := range values {
	// 	fmt.Println(v)
	// }

	output := make(chan string)
	go parser.Scan(output)

	for v := range output {
		fmt.Println(v)
	}
}
