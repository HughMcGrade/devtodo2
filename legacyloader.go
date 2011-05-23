/*
  Copyright 2011 Alec Thomas

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

// Loads legacy devtodo XML files.

package main

import (
	"io"
	"strconv"
	"strings"
	"time"
	"xml"
)

type xmlNote struct {
	Priority string "attr"
	Time string "attr"
	Done string "attr"
	Text string "chardata"
	Note []xmlNote
}

type xmlTodo struct {
	Title string
	Note []xmlNote
}

func parseXmlNote(parent Task, from []xmlNote) {
	if from == nil {
		return
	}
	for _, note := range from {
		text := strings.TrimSpace(note.Text)
		priority := PriorityFromString(note.Priority)
		task := parent.AddTask(text, priority)
		created, _ := strconv.Atoi64(note.Time)
		completed, _ := strconv.Atoi64(note.Done)
		task.SetCreationTime(time.SecondsToUTC(created))
		if completed != 0 {
			task.SetCompletionTime(time.SecondsToUTC(completed))
		}
		parseXmlNote(task, note.Note)
	}
}

func LoadLegacyTaskList(reader io.Reader) TaskList {
	var todoXml xmlTodo
	xml.Unmarshal(reader, &todoXml)

	tasks := NewTaskList()
	tasks.SetText(strings.TrimSpace(todoXml.Title))
	parseXmlNote(tasks, todoXml.Note)
	return tasks
}