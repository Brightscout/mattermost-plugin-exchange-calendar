package command

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	"time"
	//"net/http"

	"github.com/mattermost/mattermost-server/v6/model"
)

func (c *Command) setcustomstatus(parameters ...string) (string, bool, error) {
	fmt.Print("Inside setcustom status", c.Args.TriggerId)

	/*emojiList,err:=c.API.GetEmojiList("name",1,10);
	fmt.Print("emojiList",emojiList)
	fmt.Print("",)
	var mainOptions []*model.PostActionOptions
	for i, e := range emojiList {
		var options *model.PostActionOptions
		options.Value="emoji"+e.Name+string(rune(i))
		options.Text="emoji"+e.Name+string(rune(i))

		mainOptions=append(mainOptions, options)
	}
	fmt.Print("mainoptions",mainOptions)
	_ = model.OpenDialogRequest{
		TriggerId: c.Args.TriggerId,
		URL:       "http://localhost:8065/actions/dialog",
		Dialog: model.Dialog{
			Title: "Set Custom Status",
			Elements: []model.DialogElement{{
				DisplayName: "Custom Status",
				Name:        "title",
				Type:        "text",
			},/*{
				DisplayName: "Custom Emoji",
				Name: "options",
				Type: "select",
				Options: mainOptions,
			  },
		},
		},
	}
	/*var jsonData string = `  {
		"trigger_id": {{.ArgsID}},
		"url": "http://localhost:8065/store",
		"dialog": {
			"title": "test dialog",
			"introduction_text": "test introduction",
			"callback_id": "bfjsbfskf",
			"elements": [  {
							  "display_name": "Email",
							  "name": "email",
							  "type": "text",
							  "subtype": "email",
							  "placeholder": "placeholder@example.com"
						   }
						 ]
		}
	}`

	substitute := Data{c.Args.TriggerId}

	tmpl, _ := template.New("jsonData").Parse(jsonData)
	responseBody := &bytes.Buffer{}
	_ = tmpl.Execute(responseBody, substitute)

	fmt.Print("inside response\n",responseBody)
	resp, err := http.Post("http://localhost:8065/api/v4/actions/dialogs/open", "application/json", responseBody)
	   if err != nil {
		  log.Fatalf("An Error Occured %v", err)
		  fmt.Print("inside error");
	   }
	   defer resp.Body.Close()
	   //Read the response body
		  body, err := ioutil.ReadAll(resp.Body)
		  if err != nil {
			 log.Fatalln(err)
		  }
		  sb := string(body)
		  log.Print(sb)
		  fmt.Print("inside response body ioutil\n",sb);


	   fmt.Print("inside response body\n",resp.Body);
	   fmt.Print("inside response status\n",resp.Status);*/

	//by, _ := json.Marshal(request)
	/*req, _ := http.NewRequest(http.MethodPost, "http://localhost:8065/api/v4/actions/dialogs/open", bytes.NewReader(by))
	   req.Header.Add("Content-Type", "application/json")
	   resp, _ := http.DefaultClient.Do(req)
	   fmt.Printf("%#v\n", resp)
	   if err != nil {
		   fmt.Printf("%#v\n", err)
	   }*/

	customStatus := &model.CustomStatus{
		Emoji:     "calendar",
		Text:      "in a meeting",
		Duration:  "date_and_time",
		ExpiresAt: time.Now().UTC().Add(1),
	}
	fmt.Println("customStatus")
	fmt.Print(customStatus)
	fmt.Print("expiresAt", time.Now().UTC().Add(1))
	err := c.API.UpdateUserCustomStatus(c.Args.UserId, customStatus)
	if err != nil {
		fmt.Println("inside Updateusercustom Status")
	}
	return "successful", true, nil
}
