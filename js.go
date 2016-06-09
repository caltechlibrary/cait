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
	// "io/ioutil"
	"log"
	// "net/http"
	// "os"
	// "path/filepath"
	"strconv"
	// "strings"

	// 3rd Party Packages
	"github.com/robertkrimen/otto"

	// Caltech Library Packages
	"github.com/caltechlibrary/ostdlib"
)

// AddHelp adds cait API help for the JavaScript Repl
func (api *ArchivesSpaceAPI) AddHelp(js *ostdlib.JavaScriptVM) {
	js.SetHelp("api", "login", []string{}, "Using environment variables authenticate with the ArchivesSpace REST API")
	js.SetHelp("api", "logout", []string{}, "Logout of the ArchivesSpace REST API")
	js.SetHelp("api", "createRepostiory", []string{"repo object"}, "Create an repository based on the respoistory describe by the repo object")
	js.SetHelp("api", "getRepository", []string{"repoID int"}, "Get a repository object using the repository id")
	js.SetHelp("api", "updateRepository", []string{"repo object"}, "Update the repository definition using the repo object")
	js.SetHelp("api", "deleteRepository", []string{"repo object"}, "Delete the repository definition described by repo object")
	js.SetHelp("api", "listRepositories", []string{}, "list repositories ids")
	js.SetHelp("api", "createAgent", []string{"agent_type string, agent object"}, "Create an agent of agent_type using agent object. Agent types are person, software")
	js.SetHelp("api", "getAgent", []string{"agent_type string, agent_id int"}, "return an agent object using agent type and id.")
	js.SetHelp("api", "updateAgent", []string{"agent object"}, "Update an agent in ArchivesSpace from an agent object")
	js.SetHelp("api", "deleteAgent", []string{"agent object"}, "delete an agent in ArchivesSpace using agent object")
	js.SetHelp("api", "listAgents", []string{"agent_type"}, "List agents by type, returns agent ids as an array")
	js.SetHelp("api", "createAccession", []string{"repo_id int", "accession object"}, "Create an accession in the repository in ArchivesSpace from a repository id and accession object")
	js.SetHelp("api", "getAccession", []string{"repo_id int", "accession_id int"}, "Get an accession object from ArchivesSpace by repo id and accesison id")
	js.SetHelp("api", "updateAccession", []string{"accession object"}, "Update an accession from an accession object")
	js.SetHelp("api", "deleteAccession", []string{"accession object"}, "Delete an accession using an accession object")
	js.SetHelp("api", "listAccessions", []string{"repo_id int"}, "list accessions ids in repository with matching repo id")
	js.SetHelp("api", "createSubject", []string{"subject object"}, "create a subject from a subject object")
	js.SetHelp("api", "getSubject", []string{"subject_id int"}, "get an subject by subject id")
	js.SetHelp("api", "updateSubject", []string{"subject object"}, "Update a subject in ArchivesSpace from a subject object")
	js.SetHelp("api", "deleteSubject", []string{"subject object"}, "Delete a subject in ArchivesSpace from a subject object")
	js.SetHelp("api", "listSubjects", []string{}, "return a list of subject ids")
	js.SetHelp("api", "createDigitalObject", []string{"repo_id int", "digital_object object"}, "Create a digital object in the repository indicated by id")
	js.SetHelp("api", "getDigitalObject", []string{"repo_id int", "object_id int"}, "get a digital object form ArchivesSpace by repository id and object id")
	js.SetHelp("api", "updateDigitalObject", []string{"digital_object object"}, "Update a digital object in ArchivesSpace from digital_object")
	js.SetHelp("api", "deleteDigitalObject", []string{"digital_object object"}, "Delete a digital object form ArchivesSpace using digital_object")
	js.SetHelp("api", "listDigitalObjects", []string{"repo_id int"}, "get a list of digital object ids from repository with matching id")
	js.SetHelp("api", "createResource", []string{"repo_id int", "resource object"}, "Create a resource in a repository")
	js.SetHelp("api", "updateResource", []string{"resource object"}, "Update a resource")
	js.SetHelp("api", "deleteResource", []string{"resource object"}, "Delete a resource from a repository")
	js.SetHelp("api", "getResource", []string{"repo_id int", "object_id int"}, "Get a resoource from a repository")
	js.SetHelp("api", "listResources", []string{"repo_id int"}, "List resource ids in a repository")
}

// AddExtensions add cait API to a JavaScript environment along with the
func (api *ArchivesSpaceAPI) AddExtensions(js *ostdlib.JavaScriptVM) *otto.Otto {
	vm := js.VM
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
		err = ostdlib.ToStruct(call.Argument(0), &repo)
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
		err = ostdlib.ToStruct(call.Argument(0), &repo)
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
		err = ostdlib.ToStruct(call.Argument(0), &repo)
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
		err = ostdlib.ToStruct(call.Argument(1), &agent)
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
		err = ostdlib.ToStruct(call.Argument(0), &agent)
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
		err = ostdlib.ToStruct(call.Argument(0), &agent)
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
		err = ostdlib.ToStruct(call.Argument(1), &accession)
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
		err = ostdlib.ToStruct(call.Argument(0), accession)
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
		err = ostdlib.ToStruct(call.Argument(0), accession)
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
		err = ostdlib.ToStruct(call.Argument(0), subject)
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
		err = ostdlib.ToStruct(call.Argument(0), subject)
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
		err = ostdlib.ToStruct(call.Argument(0), subject)
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
		err = ostdlib.ToStruct(call.Argument(1), digitalObject)
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
		err = ostdlib.ToStruct(call.Argument(0), digitalObject)
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
		err = ostdlib.ToStruct(call.Argument(0), digitalObject)
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

	//api.createResource(repo_id int, resource object) create a resource in a repository
	apiObj.Set("createResource", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 2 {
			return errorObject(obj, fmt.Sprintf("api.createResource(repo_id, object), expects two arguments, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s := call.Argument(0).String()
		repoID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createResource(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		resourceObject := new(Resource)
		err = ostdlib.ToStruct(call.Argument(1), resourceObject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createResource(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		response, err := api.CreateResource(repoID, resourceObject)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.createResource(%s, %s) response error %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.updateResource(object) update a resource
	apiObj.Set("updateResource", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.updateResource(object), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		updateResource := new(Resource)
		err = ostdlib.ToStruct(call.Argument(0), updateResource)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateResource(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.UpdateResource(updateResource)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.updateResource(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.deleteResource(object) delete a digital object
	apiObj.Set("deleteResource", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.deleteResource(object), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		deleteResource := new(Resource)
		err = ostdlib.ToStruct(call.Argument(0), deleteResource)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteResource(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.DeleteResource(deleteResource)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.deleteResource(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.getResource(repo_id, object_id) get a digital object by repo_id and object_id
	apiObj.Set("getResource", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 2 {
			return errorObject(obj, fmt.Sprintf("api.getResource(repo_id, object_id), expects two arguments, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s := call.Argument(0).String()
		repoID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getResource(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		s = call.Argument(1).String()
		resourceID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getResource(%s, %s), arg error, %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		response, err := api.GetResource(repoID, resourceID)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.getResource(%s, %s) response error %s, %s", call.Argument(0).String(), call.Argument(1).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	// api.listResources(repo_id) list a digital objects by repo_id
	apiObj.Set("listResources", func(call otto.FunctionCall) otto.Value {
		obj, err := vm.Object(`({})`)
		if len(call.ArgumentList) != 1 {
			return errorObject(obj, fmt.Sprintf("api.listResources(repo_id), expects one argument, got %d, %s", len(call.ArgumentList), call.CallerLocation()))
		}
		s := call.Argument(0).String()
		repoID, err := strconv.Atoi(s)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.listResources(%s), arg error, %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		response, err := api.ListResources(repoID)
		if err != nil {
			return errorObject(obj, fmt.Sprintf("api.listResources(%s) response error %s, %s", call.Argument(0).String(), call.CallerLocation(), err))
		}
		return responseObject(response)
	})

	return vm
}
