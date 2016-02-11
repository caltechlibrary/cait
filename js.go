//
// Package cait is a collection of structures and functions
// for interacting with ArchivesSpace's REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package cait

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/rsdoiel/otto"
)

// NewJavaScript creates a *otto.Otto (JavaScript VM) with functions added to integrate
// the internal cait API.
func NewJavaScript(api *ArchivesSpaceAPI, jsArgs []string) *otto.Otto {
	vm := otto.New()

	errorObject := func(obj *otto.Object, msg string) otto.Value {
		if obj == nil {
			obj, _ = vm.Object(`({})`)
		}
		log.Println(msg)
		obj.Set("status", "error")
		obj.Set("error", msg)
		return obj.Value()
	}

	responseObject := func(data interface{}) otto.Value {
		src, _ := json.Marshal(data)
		obj, _ := vm.Object(fmt.Sprintf(`(%s)`, src))
		return obj.Value()
	}

	osObj, _ := vm.Object(`os = {}`)

	// os.args() returns an array of command line args
	osObj.Set("args", func(call otto.FunctionCall) otto.Value {
		results, _ := vm.ToValue(jsArgs)
		return results
	})

	// os.exit()
	osObj.Set("exit", func(call otto.FunctionCall) otto.Value {
		exitCode := 0
		if len(call.ArgumentList) == 1 {
			s := call.Argument(0).String()
			exitCode, _ = strconv.Atoi(s)
		}
		os.Exit(exitCode)
		return responseObject(exitCode)
	})

	// os.getEnv(env_varname) returns empty string or the value found as a string
	osObj.Set("getEnv", func(call otto.FunctionCall) otto.Value {
		envvar := call.Argument(0).String()
		result, err := vm.ToValue(os.Getenv(envvar))
		if err != nil {
			return errorObject(nil, fmt.Sprintf("os.getEnv(%q) %s, %s", call.CallerLocation(), envvar, err))
		}
		return result
	})

	httpObj, _ := vm.Object(`http = {}`)

	//HttpGet(uri, headers) returns contents recieved (if any)
	httpObj.Set("get", func(call otto.FunctionCall) otto.Value {
		//FIXME: Need to optional argument of an array of headers,
		// [{"Content-Type":"application/json"},{"X-ArchivesSpaceSession":"..."}]
		var headers []map[string]string

		uri := call.Argument(0).String()
		if len(call.ArgumentList) > 1 {
			rawObjs, err := call.Argument(1).Export()
			if err != nil {
				return errorObject(nil, fmt.Sprintf("Failed to process headers, %s, %s, %s", call.CallerLocation(), uri, err))
			}
			src, _ := json.Marshal(rawObjs)
			err = json.Unmarshal(src, &headers)
			if err != nil {
				return errorObject(nil, fmt.Sprintf("Failed to translate headers, %s, %s, %s", call.CallerLocation(), uri, err))
			}
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("Can't create a GET request for %s, %s, %s", uri, call.CallerLocation(), err))
		}
		for _, header := range headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("Can't connect to %s, %s, %s", uri, call.CallerLocation(), err))
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("Can't read response %s, %s, %s", uri, call.CallerLocation(), err))
		}
		return responseObject(content)
	})

	// HttpPost(uri, headers, payload) returns contents recieved (if any)
	httpObj.Set("post", func(call otto.FunctionCall) otto.Value {
		var headers []map[string]string

		uri := call.Argument(0).String()
		mimeType := call.Argument(1).String()
		payload := call.Argument(2).String()
		buf := strings.NewReader(payload)
		// Process any additional headers past to HttpPost()
		if len(call.ArgumentList) > 2 {
			rawObjs, err := call.Argument(3).Export()
			if err != nil {
				return errorObject(nil, fmt.Sprintf("Failed to process headers for %s, %s, %s", uri, call.CallerLocation(), err))
			}
			src, _ := json.Marshal(rawObjs)
			err = json.Unmarshal(src, &headers)
			if err != nil {
				return errorObject(nil, fmt.Sprintf("Failed to translate header for %s, %s, %s", uri, call.CallerLocation(), err))
			}
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", uri, buf)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("Can't create a POST request for %s, %s, %s", uri, call.CallerLocation(), err))
		}
		req.Header.Set("Content-Type", mimeType)
		for _, header := range headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("Can't connect to %s, %s, %s", uri, call.CallerLocation(), err))
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("Can't read response %s, %s, %s", uri, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(fmt.Sprintf("%s", content))
		if err != nil {
			return errorObject(nil, fmt.Sprintf("HttpGet(%q) error, %s, %s", uri, call.CallerLocation(), err))
		}
		return result
	})

	apiObj, _ := vm.Object(`api = {}`)
	// api.login() error if one occurred
	// logs into the ArchivesSpace API based on environment variables
	apiObj.Set("login", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		err = api.Login()
		if err != nil {
			return errorObject(obj, fmt.Sprintf("Login() failed, %s, %s", call.CallerLocation(), err))
		}
		obj.Set("isAuth", true)
		return obj.Value()
	})

	// logout() error if one occurred
	// logs into the ArchivesSpace API based on environment variables
	apiObj.Set("logout", func(call otto.FunctionCall) otto.Value {
		obj, _ := vm.Object(`({})`)
		err := api.Logout()
		if err != nil {
			return errorObject(obj, fmt.Sprintf("Logout() failed, %s, %s", call.CallerLocation(), err))
		}
		obj.Set("isAuth", false)
		return obj.Value()
	})

	// api.createRepository(repo) returns object of repository or error
	apiObj.Set("createRepository", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.createRepository(repository), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		repo := new(Repository)
		err = call.Argument(0).ToStruct(&repo)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createRepository() arg error %s, %s", call.CallerLocation(), err))
		}
		response, err := api.CreateRepository(repo)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createRepository() response error %s, %s", call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.getRepository(repoID) returns object of repository or error
	apiObj.Set("getRepository", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.getRepository(repoID), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s, err := call.Argument(0).ToString()
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getRepository() arg error %s, %s", call.CallerLocation(), err))
		}
		repoID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getRepository(%q) id conversion error, %s, %s", s, call.CallerLocation(), err))
		}
		response, err := api.GetRepository(repoID)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getRepository(%d) response error %s, %s", repoID, call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.updateRepository(repo) returns object of repository or error
	apiObj.Set("updateRepository", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.updateRepository(repo), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		repo := new(Repository)
		err = call.Argument(0).ToStruct(&repo)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateRepository(repo) arg error %s, %s", call.CallerLocation(), err))
		}

		response, err := api.UpdateRepository(repo)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateRepository(repo) failed to update, %s, %s", call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.deleteRepository(repo) deletes a repository and returns a response object.
	apiObj.Set("deleteRepository", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.deleteRepository(repo), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		repo := new(Repository)
		err = call.Argument(0).ToStruct(&repo)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteRepository(repo) arg error %s, %s", call.CallerLocation(), err))
		}

		response, err := api.DeleteRepository(repo)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteRepository(repo) failed to update, %s, %s", call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.listRepositories() returns a list of repository ids.
	apiObj.Set("listRepositories", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 0 {
			return errorObject(obj, fmt.Sprintf("api.listRepository(), expects zero argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}

		response, err := api.ListRepositoryIDs()
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.listRepositories() failed to update, %s, %s", call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.createAgent(agent_type, agent) creates a new agent of agent_type
	apiObj.Set("createAgent", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 2 {
			return errorObject(obj, fmt.Sprintf("api.createAgent(agent_type, agent), expects two arguments, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		agentType := call.Argument(0).String()
		agent := new(Agent)
		err = call.Argument(1).ToStruct(&agent)
		if err != nil || agentType == "" {
			return errorObject(obj, fmt.Sprintf("api.createAgent(agent_type, agent) arg error %s, %s", call.CallerLocation(), err))
		}
		response, err := api.CreateAgent(agentType, agent)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createAgent(agent_type, agent) response error %s, %s", call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.getAgent(agent_type, agent_id) get an agent of agent_type
	apiObj.Set("getAgent", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 2 {
			return errorObject(obj, fmt.Sprintf("api.getAgent(agent_type, agent_id), expects two arguments, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		agentType := call.Argument(0).String()
		agentID := 0
		s := call.Argument(1).String()
		agentID, err = strconv.Atoi(s)
		if err != nil || agentType == "" || agentID == 0 {
			return errorObject(obj, fmt.Sprintf("api.getAgent(%q, %d) arg error %s, %s", agentType, agentID, call.CallerLocation(), err))
		}
		response, err := api.GetAgent(agentType, agentID)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getAgent(%q, %d) response error %s, %s", agentType, agentID, call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.updateAgent(agent_type, agent_id) update an agent of agent_type
	apiObj.Set("updateAgent", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.updateAgent(agent), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		agent := new(Agent)
		err = call.Argument(0).ToStruct(&agent)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateAgent(%q) arg error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.UpdateAgent(agent)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateAgent(%q) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.updateAgent(agent) update an agent
	apiObj.Set("updateAgent", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.updateAgent(agent), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		agent := new(Agent)
		err = call.Argument(0).ToStruct(&agent)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateAgent(%q) arg error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.UpdateAgent(agent)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateAgent(%q) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.deleteAgent(agent) delete an agent
	apiObj.Set("deleteAgent", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.deleteAgent(agent), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		agent := new(Agent)
		err = call.Argument(0).ToStruct(&agent)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteAgent(%q) arg error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.DeleteAgent(agent)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteAgent(%q) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.listAgents(agent_type) list agent ids
	apiObj.Set("listAgents", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.listAgents(agent_type), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		agentType := call.Argument(0).String()
		response, err := api.ListAgents(agentType)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.listAgents(%q) response error %s, %s", agentType, call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.createAccession(repo_id, accession) create an accession
	apiObj.Set("createAccession", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 2 {
			return errorObject(obj, fmt.Sprintf("api.createAccession(repo_id, accession), expects two arguments, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s := call.Argument(0).String()
		repoID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createAccession(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		accession := new(Accession)
		err = call.Argument(1).ToStruct(&accession)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createAccession(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		response, err := api.CreateAccession(repoID, accession)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createAccession(%s, %s) response error %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.getAccession(repo_id, accession_id) get an accession
	apiObj.Set("getAccession", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 2 {
			return errorObject(obj, fmt.Sprintf("api.getAccession(repo_id, accession_id), expects two arguments, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s := call.Argument(0).String()
		repoID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getAccession(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		s = call.Argument(1).String()
		accessionID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getAccession(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		response, err := api.GetAccession(repoID, accessionID)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createAccession(%s, %s) response error %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.updateAccession(accession) update an accession
	apiObj.Set("updateAccession", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.updateAccession(accession), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		accession := new(Accession)
		err = call.Argument(0).ToStruct(accession)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateAccession(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.UpdateAccession(accession)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateAccession(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.deleteAccession(accession) delete an accession
	apiObj.Set("deleteAccession", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.deleteAccession(accession), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		accession := new(Accession)
		err = call.Argument(0).ToStruct(accession)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteAccession(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.DeleteAccession(accession)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateAccession(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.listAccessions(repo_id) list accession ids
	apiObj.Set("listAccessions", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.listAccessions(repo_id), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s := call.Argument(0).String()
		repoID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.listAccession(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.ListAccessions(repoID)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.listAccessions(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.createSubject(subject) create a subject
	apiObj.Set("createSubject", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.createSubject(subject), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		subject := new(Subject)
		err = call.Argument(0).ToStruct(subject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createSubject(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.CreateSubject(subject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createSubject(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.getSubject(subject_id) get a subject
	apiObj.Set("getSubject", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.getSubject(subject_id), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s := call.Argument(0).String()
		subjectID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getSubject(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.GetSubject(subjectID)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getSubject(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.updateSubject(subject) update a subject
	apiObj.Set("updateSubject", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.updateSubject(subject), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		subject := new(Subject)
		err = call.Argument(0).ToStruct(subject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateSubject(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.UpdateSubject(subject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getSubject(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.deleteSubject(subject) delete a subject
	apiObj.Set("deleteSubject", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.deleteSubject(subject), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		subject := new(Subject)
		err = call.Argument(0).ToStruct(subject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteSubject(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.UpdateSubject(subject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteSubject(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.listSubjects() list subject ids
	apiObj.Set("listSubjects", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 0 {
			return errorObject(obj, fmt.Sprintf("api.listSubjects(), expects zero argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		response, err := api.ListSubjects()
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.listSubjects() response error %s, %s", call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.createDigitalObject(repo_id, object) create a digital object
	apiObj.Set("createDigitalObject", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 2 {
			return errorObject(obj, fmt.Sprintf("api.getDigitalObject(repo_id, object), expects two arguments, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s := call.Argument(0).String()
		repoID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createDigitalObject(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		digitalObject := new(DigitalObject)
		err = call.Argument(1).ToStruct(digitalObject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createDigitalObject(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		response, err := api.CreateDigitalObject(repoID, digitalObject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createDigitalObject(%s, %s) response error %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.getDigitalObject(repo_id, object_id) get a digital object by repo_id and object_id
	apiObj.Set("getDigitalObject", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 2 {
			return errorObject(obj, fmt.Sprintf("api.getDigitalObject(repo_id, object_id), expects two arguments, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s := call.Argument(0).String()
		repoID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getDigitalObject(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		s = call.Argument(1).String()
		objectID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getDigitalObject(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		response, err := api.GetDigitalObject(repoID, objectID)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getDigitalObject(%s, %s) response error %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.updateDigitalObject(object) update a digital object
	apiObj.Set("updateDigitalObject", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.updateDigitalObject(object), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		digitalObject := new(DigitalObject)
		err = call.Argument(0).ToStruct(digitalObject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateDigitalObject(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.UpdateDigitalObject(digitalObject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateDigitalObject(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.deleteDigitalObject(object) delete a digital object
	apiObj.Set("deleteDigitalObject", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.deleteDigitalObject(object), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		digitalObject := new(DigitalObject)
		err = call.Argument(0).ToStruct(digitalObject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteDigitalObject(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.DeleteDigitalObject(digitalObject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteDigitalObject(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.listDigitalObjects(repo_id) list a digital objects by repo_id
	apiObj.Set("listDigitalObjects", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.listDigitalObjects(repo_id), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s := call.Argument(0).String()
		repoID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.listDigitalObjects(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.ListDigitalObjects(repoID)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.listDigitalObjects(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	//
	// Add Polyfills, FIXME: these need to be implemented in Otto...
	//
	vm.Eval(`if (!String.prototype.repeat) {
	  String.prototype.repeat = function(count) {
	    'use strict';
	    if (this == null) {
	      throw new TypeError('can\'t convert ' + this + ' to object');
	    }
	    var str = '' + this;
	    count = +count;
	    if (count != count) {
	      count = 0;
	    }
	    if (count < 0) {
	      throw new RangeError('repeat count must be non-negative');
	    }
	    if (count == Infinity) {
	      throw new RangeError('repeat count must be less than infinity');
	    }
	    count = Math.floor(count);
	    if (str.length == 0 || count == 0) {
	      return '';
	    }
	    // Ensuring count is a 31-bit integer allows us to heavily optimize the
	    // main part. But anyway, most current (August 2014) browsers can't handle
	    // strings 1 << 28 chars or longer, so:
	    if (str.length * count >= 1 << 28) {
	      throw new RangeError('repeat count must not overflow maximum string size');
	    }
	    var rpt = '';
	    for (;;) {
	      if ((count & 1) == 1) {
	        rpt += str;
	      }
	      count >>>= 1;
	      if (count == 0) {
	        break;
	      }
	      str += str;
	    }
	    // Could we try:
	    // return Array(count + 1).join(this);
	    return rpt;
	  }
	}
`)

	return vm
}
