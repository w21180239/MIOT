package server

import (
	"net/http"

	"gitlab.com/LICOTEK/DuerOS/duerosprotocol"
)

func handlerProtocol() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("Header：")
		//if len(r.Header) > 0 {
		//	for k,v := range r.Header {
		//		fmt.Printf("%s=%s\n", k, v[0])
		//	}
		//}

		//fmt.Println("Body：")
		//buf := new(bytes.Buffer)
		//buf.ReadFrom(r.Body)
		//fmt.Println(buf.String())

		processor := duerosprotocol.NewProcessor(r)
		result := processor.Process()
		formatter.JSON(w, http.StatusOK, result)
	}

}
