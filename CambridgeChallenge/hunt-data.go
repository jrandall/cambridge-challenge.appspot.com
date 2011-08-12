/*
 hunt-data

    Copyright 2011 Joshua C. Randall <jcrandall@alum.mit.edu>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package CambridgeChallenge

import (
    "os"
    "json"
    "fmt"
)

type Clue struct {
     Name string
     Prompt string
     Answer string
}

type State struct {
     Name string
     NextState string
     Clues []Clue
}

type Hunt struct {
     Name string
     Date string
     EnterState string
     States []State
}


var h Hunt

func init() {
     // CreateHuntDataTemplate("hunt-template.json")
     h = LoadHuntData("hunt.json")
}

func CreateHuntDataTemplate(saveFileName string) {
     huntFileWriter, err := os.Create(saveFileName)
     	if err != nil {
	  // error opening hunt file
       	  panic(fmt.Sprintf("Error creating file %v", err))
     	} else {
	  // template hunt data file created and opened for writing
	  defer huntFileWriter.Close()
	  encoder := json.NewEncoder(huntFileWriter)
	  dummyClueA1 := Clue{Name:"NameA1",Prompt:"PromptA1",Answer:"AnswerA1"}
	  dummyClueA2 := Clue{Name:"NameA2",Prompt:"PromptA2",Answer:"AnswerA2"}
	  dummyStateA := State{Name:"NameA",Clues:[]Clue{dummyClueA1,dummyClueA2}}
	  dummyClueB1 := Clue{Name:"NameB1",Prompt:"PromptB1",Answer:"AnswerB1"}
	  dummyClueB2 := Clue{Name:"NameB2",Prompt:"PromptB2",Answer:"AnswerB2"}
	  dummyStateB := State{Name:"NameB",Clues:[]Clue{dummyClueB1,dummyClueB2}}
	  huntTemplate := Hunt{Name:"HuntName",Date:"HuntDate",States:[]State{dummyStateA,dummyStateB}}
	  fmt.Println("json encoding hunt", huntTemplate)
	  err := encoder.Encode(&huntTemplate)
	  if err != nil {
	    panic(fmt.Sprintf("Error encoding JSON %v", err))
	  }
	}
	return
}

func LoadHuntData(huntFileName string) (h Hunt) {
     	huntFileReader, err := os.Open(huntFileName)
     	if err != nil {
	  // error opening hunt file
       	  panic(fmt.Sprintf("Error %v", err))
     	} else {
	  // hunt file opened
	  defer huntFileReader.Close()
          decoder := json.NewDecoder(huntFileReader)
	  err := decoder.Decode(&h)
	  if err != nil {
	    panic(fmt.Sprintf("Error decoding JSON %v", err))
	  }
     	}
     	return h
}

