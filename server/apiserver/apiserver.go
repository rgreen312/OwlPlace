package apiserver

import (
	"os"
	"fmt"
	"net/http"
	"github.com/rgreen312/owlplace/server/consensus"
	"html/template"

)


type ApiServer struct {
	sendc chan consensus.BackendMessage
	recvc chan consensus.ConsensusMessage
}

func NewApiServer(send_channel chan consensus.BackendMessage, recv_channel chan consensus.ConsensusMessage) *ApiServer {
	return &ApiServer{
		sendc: send_channel,
		recvc: recv_channel,
	}
}

func (api *ApiServer) ListenAndServe(){
	http.HandleFunc("/headers", api.Headers)
	http.HandleFunc("/get_image", api.GetImage)
	http.HandleFunc("/update_pixel", api.UpdatePixel)
	http.ListenAndServe(":3000", nil)
}


func (api *ApiServer) GetImage(w http.ResponseWriter, req *http.Request) {
	// Debug message
	fmt.Fprintf(os.Stdout, "Getting Image From Raft\n")
	// Construct the message
	m := consensus.BackendMessage{ Type: consensus.GET_IMAGE }
	// Send a message through the channel
	api.sendc <- m
	var ImageTemplate string = `<!DOCTYPE html>
								<html lang="en"><head></head>
								<body><img src="data:image/jpg;base64,{{.Image}}"></body>`
	image_msg := <- api.recvc

	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		fmt.Fprintf(os.Stdout, "Unable to parse image template.\n")
	} else {
		data := map[string]interface{}{"Image": image_msg.Data}
		if err = tmpl.Execute(w, data); err != nil {
			fmt.Fprintf(os.Stdout, "Unable to execute template.\n")
		}
	}

}

func (api *ApiServer) UpdatePixel(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(os.Stdout, "Not Implemented.\n")
}


func (api *ApiServer) Headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
