package grabber

import (
	"encoding/json"
)

type FirstTask struct {
	InputFlowData struct {
		FlowContext struct {
			DebugOverrides map[string]interface{} `json:"debug_overrides"`
			StartLocation  struct {
				Location string `json:"location"`
			} `json:"start_location"`
		} `json:"flow_context"`
	} `json:"input_flow_data"`
	SubtaskVersions map[string]int `json:"subtask_versions"`
}

func GenerateFirstTaskString() string {
	task := FirstTask{
		InputFlowData: struct {
			FlowContext struct {
				DebugOverrides map[string]interface{} `json:"debug_overrides"`
				StartLocation  struct {
					Location string `json:"location"`
				} `json:"start_location"`
			} `json:"flow_context"`
		}{
			FlowContext: struct {
				DebugOverrides map[string]interface{} `json:"debug_overrides"`
				StartLocation  struct {
					Location string `json:"location"`
				} `json:"start_location"`
			}{
				DebugOverrides: make(map[string]interface{}),
				StartLocation: struct {
					Location string `json:"location"`
				}{
					Location: "manual_link",
				},
			},
		},
		SubtaskVersions: map[string]int{
			"action_list": 2,
			// Добавьте остальные версии подзадач здесь...
		},
	}

	jsonData, _ := json.Marshal(task)
	return string(jsonData)
}

type SecondTask struct {
	FlowToken     string `json:"flow_token"`
	SubtaskInputs []struct {
		SubtaskID         string `json:"subtask_id"`
		JsInstrumentation struct {
			Response string `json:"response"`
			Link     string `json:"link"`
		} `json:"js_instrumentation"`
	} `json:"subtask_inputs"`
}

func GenerateSecondTaskString(flowToken string) string {
	task := SecondTask{
		FlowToken: flowToken,
		SubtaskInputs: []struct {
			SubtaskID         string `json:"subtask_id"`
			JsInstrumentation struct {
				Response string `json:"response"`
				Link     string `json:"link"`
			} `json:"js_instrumentation"`
		}{
			{
				SubtaskID: "LoginJsInstrumentationSubtask",
				JsInstrumentation: struct {
					Response string `json:"response"`
					Link     string `json:"link"`
				}{
					Response: `{"rf":{"a41eea33ff45860584e0a227cca4d3b195dc4c024f64f5e07f7f08c1f0045773":-1,"a3fbe732c22c0e16bd3424dd41be249e363d8bf64e46c5349aba728a1d99001b":95,"aff908271b9a8bc0140fcc1a221cb3bd6aafaab0ee4e8e246791a17eeda505da":234,"dcbb29c850b8d22588a377fafa2d2f742c38c6ecf20874c8b14e6b62e0473f77":180},"s":"YrwdUofQdV3MssGow7mi0ew0No0tgY57AdhypjW5h75VBvTMcWT5ClJmAXJO33usTnzi5-PxsykcIZ5IgJCnSJB7BUf2FLq946udjqKjT_CglZluS-if5E8nm9DPuHuLwF2E76XttS9lpxVST_eFfnl96EtadnGTxuLa471S0_1LcXxdrqXQpqkObHYNchGymzxTT1zc-r9KkBo800CXo5S32s0pE9eRZuBjOPGiTK0QHLPGZsrhOLTHUrs9yhJQPKfUAGz7zEIZe0sZVHggFZWpsuc9ozoD-Omxfpw3JZqiZ-JovhnPyWICzz4S58SXYbQBmfH_btRkvXM8Ad0uVgAAAYuP1-lw"}`, Link: "next_link",
				},
			},
		},
	}
	jsonData, _ := json.Marshal(task)
	return string(jsonData)
}

type ThirdTask struct {
	FlowToken     string `json:"flow_token"`
	SubtaskInputs []struct {
		SubtaskID    string `json:"subtask_id"`
		SettingsList struct {
			SettingResponses []struct {
				Key          string `json:"key"`
				ResponseData struct {
					TextData struct {
						Result string `json:"result"`
					} `json:"text_data"`
				} `json:"response_data"`
			} `json:"setting_responses"`
			Link string `json:"link"`
		} `json:"settings_list"`
	} `json:"subtask_inputs"`
}

func GenerateThirdTaskString(flowToken string, twitterLogin string) string {
	task := ThirdTask{
		FlowToken: flowToken,
		SubtaskInputs: []struct {
			SubtaskID    string `json:"subtask_id"`
			SettingsList struct {
				SettingResponses []struct {
					Key          string `json:"key"`
					ResponseData struct {
						TextData struct {
							Result string `json:"result"`
						} `json:"text_data"`
					} `json:"response_data"`
				} `json:"setting_responses"`
				Link string `json:"link"`
			} `json:"settings_list"`
		}{
			{
				SubtaskID: "LoginEnterUserIdentifierSSO",
				SettingsList: struct {
					SettingResponses []struct {
						Key          string `json:"key"`
						ResponseData struct {
							TextData struct {
								Result string `json:"result"`
							} `json:"text_data"`
						} `json:"response_data"`
					} `json:"setting_responses"`
					Link string `json:"link"`
				}{
					SettingResponses: []struct {
						Key          string `json:"key"`
						ResponseData struct {
							TextData struct {
								Result string `json:"result"`
							} `json:"text_data"`
						} `json:"response_data"`
					}{
						{
							Key: "user_identifier",
							ResponseData: struct {
								TextData struct {
									Result string `json:"result"`
								} `json:"text_data"`
							}{
								TextData: struct {
									Result string `json:"result"`
								}{
									Result: twitterLogin,
								},
							},
						},
					},
					Link: "next_link",
				},
			},
		},
	}

	jsonData, _ := json.Marshal(task)
	return string(jsonData)
}

type FourthTask struct {
	FlowToken     string `json:"flow_token"`
	SubtaskInputs []struct {
		SubtaskID     string `json:"subtask_id"`
		EnterPassword struct {
			Password string `json:"password"`
			Link     string `json:"link"`
		} `json:"enter_password"`
	} `json:"subtask_inputs"`
}

func GenerateFourthTaskString(flowToken string, twitterPassword string) string {
	task := FourthTask{
		FlowToken: flowToken,
		SubtaskInputs: []struct {
			SubtaskID     string `json:"subtask_id"`
			EnterPassword struct {
				Password string `json:"password"`
				Link     string `json:"link"`
			} `json:"enter_password"`
		}{
			{
				SubtaskID: "LoginEnterPassword",
				EnterPassword: struct {
					Password string `json:"password"`
					Link     string `json:"link"`
				}{
					Password: twitterPassword,
					Link:     "next_link",
				},
			},
		},
	}

	jsonData, _ := json.Marshal(task)
	return string(jsonData)
}

type FifthTask struct {
	FlowToken     string `json:"flow_token"`
	SubtaskInputs []struct {
		SubtaskID            string `json:"subtask_id"`
		CheckLoggedInAccount struct {
			Link string `json:"link"`
		} `json:"check_logged_in_account"`
	} `json:"subtask_inputs"`
}

func GenerateFifthTaskString(flowToken string) string {
	task := FifthTask{
		FlowToken: flowToken,
		SubtaskInputs: []struct {
			SubtaskID            string `json:"subtask_id"`
			CheckLoggedInAccount struct {
				Link string `json:"link"`
			} `json:"check_logged_in_account"`
		}{
			{
				SubtaskID: "AccountDuplicationCheck",
				CheckLoggedInAccount: struct {
					Link string `json:"link"`
				}{
					Link: "AccountDuplicationCheck_false",
				},
			},
		},
	}

	jsonData, _ := json.Marshal(task)
	return string(jsonData)
}
