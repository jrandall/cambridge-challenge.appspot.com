/*
 hunt-data.go - data structures and reader/writer for hunt data

    Copyright 2011 Joshua C. Randall <jcrandall@alum.mit.edu>

    This file is part of CambridgeChallenge.

    CambridgeChallenge is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    CambridgeChallenge is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with CambridgeChallenge.  If not, see <http://www.gnu.org/licenses/>.
*/
package CambridgeChallenge

import (
    "os"
    "json"
    "fmt"
    "io"
)

type Clue struct {
     ClueName string // optional, default ""
     Prompt string   // optional, default ""
     Answer string   // optional, default ""
     AnswerType string // optional, default "regexp"
     CorrectAnswerState string // optional, default "" (stay in current state)
     IncorrectAnswerState string // optional, default "" (stay in current state)
}

type State struct {
     StateName string
     NextState string
     PreviousState string
     AllowNetMask string
     Clues []Clue
}

type Hunt struct {
     HuntName string
     HuntDate string
     EnterState string
     States map[string]*State
}


//var h Hunt

func init() {
     // CreateHuntDataTemplateFile("hunt-template.json")
     //     h = LoadHuntDataFile("hunt.json")
}

func CreateHuntDataTemplateFile(saveFileName string) {
     huntFileWriter, err := os.Create(saveFileName)
     	if err != nil {
	  // error opening hunt file
       	  panic(fmt.Sprintf("Error creating file %v", err))
     	} else {
	  // template hunt data file created and opened for writing
	  defer huntFileWriter.Close()
	  dummyClueA1 := Clue{ClueName:"NameA1",Prompt:"PromptA1",Answer:"AnswerA1"}
	  dummyClueA2 := Clue{ClueName:"NameA2",Prompt:"PromptA2",Answer:"AnswerA2"}
	  dummyStateA := State{StateName:"NameA",Clues:[]Clue{dummyClueA1,dummyClueA2}}
	  dummyClueB1 := Clue{ClueName:"NameB1",Prompt:"PromptB1",Answer:"AnswerB1"}
	  dummyClueB2 := Clue{ClueName:"NameB2",Prompt:"PromptB2",Answer:"AnswerB2"}
	  dummyStateB := State{StateName:"NameB",Clues:[]Clue{dummyClueB1,dummyClueB2}}
	  huntTemplate := Hunt{
	  	       HuntName:"HuntName",
	  	       HuntDate:"HuntDate",
		       States:map[string]*State{"NameA":&dummyStateA,"NameB":&dummyStateB}}
	  fmt.Println("json encoding hunt", huntTemplate)
	  err = EncodeHuntData(huntFileWriter, &huntTemplate)
	  if err != nil {
	    panic(fmt.Sprintf("Error encoding JSON %v", err))
	  }
	}
	return
}

func LoadHuntDataFile(huntFileName string) (h *Hunt) {
     	huntFileReader, err := os.Open(huntFileName)
     	if err != nil {
	  // error opening hunt file
       	  panic(fmt.Sprintf("Error %v", err))
     	} else {
	  // hunt file opened
	  defer huntFileReader.Close()
	  h, err = DecodeHuntData(huntFileReader)
	  if err != nil {
	    panic(fmt.Sprintf("Error decoding JSON %v", err))
	  }
	}
	return // h	  
}

func DecodeHuntData(huntJSONReader io.Reader) (h *Hunt, err os.Error) {
        decoder := json.NewDecoder(huntJSONReader)
	err = decoder.Decode(&h)
     	return // h, err
}

func EncodeHuntData(huntJSONWriter io.Writer, h *Hunt) (err os.Error) {
	  encoder := json.NewEncoder(huntJSONWriter)
	  err = encoder.Encode(h)
	  return // err
}
