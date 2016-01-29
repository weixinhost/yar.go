package main

import (
	"yar/packager/json"
	"fmt"
	"unsafe"
)

type ColorGroup struct {
	ID     int			`json:"id"`
	Name   string		`json:"name"`
	Colors []string		`json:"colors"`
}

func main(){
/*


	type ColorGroup2 struct {
		ID     int			`json:"id"`
		Name   string		`json:"name"`
		Colors string		`json:"colors"`
	}



	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}

	g := ColorGroup2{}

	json.Unmarshal(b,&g)

	os.Stdout.Write(b)

	fmt.Printf("%d %s",g.ID,g.Colors)

*/

	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}

	ret,err := json.Pack(&group);

	if err != nil {

		fmt.Printf("%s",err)

	}

	group2 := ColorGroup{}

	json.Unpack(ret,&group2)

	fmt.Printf("Pack:%s",ret)

	fmt.Printf("%d %s",group2.ID,group2.Name)

	fmt.Printf("\n\n%d\n\n",unsafe.Sizeof(group2))

}
