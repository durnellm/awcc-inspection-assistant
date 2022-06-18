package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"

	tcell "github.com/gdamore/tcell/v2"
	godotenv "github.com/joho/godotenv"
	tview "github.com/rivo/tview"
)

type Entries struct {
	Entry  []Entry `json:"data"`
	Paging Paging  `json:"paging`
}

type Paging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type Entry struct {
	Node        Node        `json:"node"`
	List_status List_status `json:"list_status"`
}

type Node struct {
	Id           int          `json:"id"`
	Title        string       `json:"title"`
	Main_picture Main_picture `json:"main_picture"`
	Media        string       `json:"media_type"`
}

type Main_picture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type List_status struct {
	Status               string `json:"status"`
	Score                int    `json:"score"`
	Num_episodes_watched int    `json:"num_episodes_watched"`
	Is_rewatching        bool   `json:"is_rewatching"`
	Updated_at           string `json:"updated_at"`
	Start_date           string `json:"start_date"`
	Finish_date          string `json:"finish_date"`
}

type CleanEntry struct {
	id        string
	title     string
	startdate string
	enddate   string
	maltype   string
}

type Filter_Data struct {
	Username     string
	StartMonth   string
	StartDay     string
	StartYear    string
	EndMonth     string
	EndDay       string
	EndYear      string
	PrevStart    bool
	PrevComp     bool
	TV           bool
	OVA          bool
	ONA          bool
	Special      bool
	Movie        bool
	Music        bool
	Unknown      bool
	anime        []Entry
	FilteredList []CleanEntry
	SortedList   []CleanEntry
}

func Date_filter(list []Entry, start, end string, prev_start, prev_comp bool) (filtered_list []CleanEntry) {
	filtered_list = []CleanEntry{}
	temp := CleanEntry{}
	if prev_comp {
		for i := 0; i < len(list); i++ {
			if list[i].List_status.Finish_date <= end {
				temp.id = strconv.Itoa(list[i].Node.Id)
				temp.title = list[i].Node.Title
				temp.startdate = list[i].List_status.Start_date
				temp.enddate = list[i].List_status.Finish_date
				temp.maltype = list[i].Node.Media
				filtered_list = append(filtered_list, temp)
			}
		}
	} else if prev_start {
		for i := 0; i < len(list); i++ {
			if list[i].List_status.Finish_date <= end && list[i].List_status.Finish_date >= start {
				temp.id = strconv.Itoa(list[i].Node.Id)
				temp.title = list[i].Node.Title
				temp.startdate = list[i].List_status.Start_date
				temp.enddate = list[i].List_status.Finish_date
				temp.maltype = list[i].Node.Media
				filtered_list = append(filtered_list, temp)
			}
		}
	} else {
		for i := 0; i < len(list); i++ {
			if list[i].List_status.Finish_date <= end && list[i].List_status.Start_date >= start {
				temp.id = strconv.Itoa(list[i].Node.Id)
				temp.title = list[i].Node.Title
				temp.startdate = list[i].List_status.Start_date
				temp.enddate = list[i].List_status.Finish_date
				temp.maltype = list[i].Node.Media
				filtered_list = append(filtered_list, temp)
			}
		}
	}
	return
}
func Type_filter(list []CleanEntry, tv, movie, ova, ona, special, music, unknown bool) (filtered_list []CleanEntry) {
	var media_types = []string{}
	if tv {
		media_types = append(media_types, "tv")
	}
	if movie {
		media_types = append(media_types, "movie")
	}
	if ova {
		media_types = append(media_types, "ova")
	}
	if ona {
		media_types = append(media_types, "ona")
	}
	if special {
		media_types = append(media_types, "special")
	}
	if music {
		media_types = append(media_types, "music")
	}
	if unknown {
		media_types = append(media_types, "unknown")
	}

	for i := 0; i < len(list); i++ {
		for _, val := range media_types {
			if val == list[i].maltype {
				filtered_list = append(filtered_list, list[i])
			}
		}
	}
	return
}

func Date_sort(list []CleanEntry) (sorted_list []CleanEntry) {
	for i := 0; i < len(list); i++ {
		sorted_list = append(sorted_list, list[i])
	}
	sort.SliceStable(sorted_list, func(i, j int) bool {
		return sorted_list[i].enddate < sorted_list[j].enddate
	})

	return
}

func Get_list(username string) (list []Entry) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	secretKey := os.Getenv("SECRET_KEY")

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.myanimelist.net/v2/users/"+username+"/animelist?status=completed&limit=1000&fields=list_status,media_type&nsfw=true", nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header = http.Header{
		"X-MAL-CLIENT-ID": {secretKey},
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	byteValue, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	var temp Entries

	err = json.Unmarshal(byteValue, &temp)

	if err != nil {
		fmt.Println(err)
	}

	next := temp.Paging.Next

	for next != "" {
		req, err := http.NewRequest("GET", next, nil)
		if err != nil {
			fmt.Println(err)
		}

		req.Header = http.Header{
			"X-MAL-CLIENT-ID": {secretKey},
		}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
		}

		if resp.Body != nil {
			defer resp.Body.Close()
		}

		byteValue, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err)
		}

		var temp2 Entries

		err = json.Unmarshal(byteValue, &temp2)

		if err != nil {
			fmt.Println(err)
		}

		for i := 0; i < len(temp2.Entry); i++ {
			temp.Entry = append(temp.Entry, temp2.Entry[i])
		}

		next = temp2.Paging.Next
	}

	list = temp.Entry

	return
}

var Once = true
var app = tview.NewApplication()
var pages = tview.NewPages()
var top = tview.NewForm()
var filter = tview.NewFlex()
var ListA = tview.NewFlex()
var ListE = tview.NewFlex()
var filtersA = tview.NewForm()
var filtersB = tview.NewForm()
var filtersZ = tview.NewForm()
var typefilters = tview.NewForm()
var listdisplayA = tview.NewTable().
	SetEvaluateAllRows(true).
	SetSelectable(true, false)
var listdisplayE = tview.NewTable().
	SetEvaluateAllRows(true).
	SetSelectable(true, false)
var Months = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

var data = Filter_Data{}

var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("Press (Enter) to continue")

var filttext = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("Filter by type (leave blank to skip)")

var untext = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetTextAlign(1)

var listinfo = tview.NewTextView().
	SetTextColor(tcell.ColorGreen)

func mainForm() *tview.Form {
	top.AddInputField("Username", "", 20, nil, func(username string) {
		data.Username = username
	})

	top.AddButton("Continue", func() {
		if data.Username != "" {
			untext.SetText(data.Username + "'s List")
			if Once {
				filterstart()
				filterend()
				typebuttons()
				filterbuttons()
			}
			Once = false
			data.anime = Get_list(data.Username)
			pages.SwitchToPage("Filter")
		}
	})

	return top
}

func filterstart() *tview.Form {
	filtersA.SetHorizontal(true)

	filtersA.AddInputField("Start Date (YYYY/MM/DD)", "", 4, nil, func(startYear string) {
		data.StartYear = startYear
	})

	filtersA.AddDropDown("", Months, 0, func(month string, index int) {
		data.StartMonth = strconv.Itoa(index + 1)
		if len(data.StartMonth) == 1 {
			data.StartMonth = "0" + data.StartMonth
		}
	})

	filtersA.AddInputField("", "", 2, nil, func(startDay string) {
		data.StartDay = startDay
		if len(data.StartDay) == 1 {
			data.StartDay = "0" + data.StartDay
		}
	})

	filtersA.AddCheckbox("Prev. Started", false, func(prevStart bool) {
		data.PrevStart = prevStart
	})
	return filtersA
}

func filterend() *tview.Form {
	filtersB.SetHorizontal(true)

	filtersB.AddInputField("End Date (YYYY/MM/DD)", "", 4, nil, func(endYear string) {
		data.EndYear = endYear
	})

	filtersB.AddDropDown("", Months, 0, func(month string, index int) {
		data.EndMonth = strconv.Itoa(index + 1)
		if len(data.EndMonth) == 1 {
			data.EndMonth = "0" + data.EndMonth
		}
	})

	filtersB.AddInputField("", "", 2, nil, func(endDay string) {
		data.EndDay = endDay
		if len(data.EndDay) == 1 {
			data.EndDay = "0" + data.EndDay
		}
	})

	filtersB.AddCheckbox("Prev. Completed", false, func(prevComp bool) {
		data.PrevComp = prevComp
	})
	return filtersB
}
func filterbuttons() *tview.Form {
	filtersZ.SetHorizontal(true)

	filtersZ.AddButton("Continue", func() {
		data.FilteredList = Date_filter(data.anime, data.StartYear+"-"+data.StartMonth+"-"+data.StartDay, data.EndYear+"-"+data.EndMonth+"-"+data.EndDay, data.PrevStart, data.PrevComp)
		if data.TV || data.Movie || data.OVA || data.ONA || data.Special || data.Music || data.Unknown {
			data.FilteredList = Type_filter(data.FilteredList, data.TV, data.Movie, data.OVA, data.ONA, data.Special, data.Music, data.Unknown)
		}
		data.SortedList = Date_sort(data.FilteredList)

		listinfo.SetText("Sort: (1) Alphabetically \t (2) End Date \t\t Control: (5) Back \t (6) Start\n" + data.Username + ": " + data.StartYear + "-" + data.StartMonth + "-" + data.StartDay + " " + data.EndYear + "-" + data.EndMonth + "-" + data.EndDay)
		listdisplayA.Clear()
		listdisplayE.Clear()
		fillTableAlph()
		pages.SwitchToPage("ListAlph")
	})
	return filtersZ
}

func typebuttons() *tview.Form {
	typefilters.SetHorizontal(true)

	typefilters.AddCheckbox("TV", false, func(tv bool) {
		data.TV = tv
	})

	typefilters.AddCheckbox("Movie", false, func(movie bool) {
		data.Movie = movie
	})

	typefilters.AddCheckbox("OVA", false, func(ova bool) {
		data.OVA = ova
	})

	typefilters.AddCheckbox("ONA", false, func(ona bool) {
		data.ONA = ona
	})

	typefilters.AddCheckbox("Special", false, func(special bool) {
		data.Special = special
	})

	typefilters.AddCheckbox("Music", false, func(music bool) {
		data.Music = music
	})

	typefilters.AddCheckbox("Unknown", false, func(unknown bool) {
		data.Unknown = unknown
	})
	return typefilters
}

func display(listdisplay tview.Table, list []CleanEntry) {
	c := 0
	for i := 0; i < len(list); i++ {
		var count = tview.NewTableCell("")
		var id = tview.NewTableCell("")
		var titles []string
		var mtype = tview.NewTableCell("")
		var start = tview.NewTableCell("")
		var end = tview.NewTableCell("")

		count.SetText(strconv.Itoa(i + 1))
		id.SetText(list[i].id)
		if len(list[i].title) > 50 {
			temp := tview.WordWrap(list[i].title, 50)
			for j := 0; j < len(temp); j++ {
				titles = append(titles, temp[j])
			}
		} else {
			titles = append(titles, list[i].title)
		}
		mtype.SetText(list[i].maltype)
		start.SetText(list[i].startdate)
		end.SetText(list[i].enddate)

		listdisplay.SetCell(i+c, 0, count)
		listdisplay.SetCell(i+c, 1, id)
		if len(titles) > 1 {
			for j := 0; j < len(titles); j++ {
				listdisplay.SetCell(i+c, 2, tview.NewTableCell(titles[j]))
				c = c + 1
			}
			c = c - 1
		} else {
			listdisplay.SetCell(i+c, 2, tview.NewTableCell(titles[0]))
		}
		listdisplay.SetCell(i+c, 3, mtype)
		listdisplay.SetCell(i+c, 4, start)
		listdisplay.SetCell(i+c, 5, end)
	}
}

func fillTableAlph() *tview.Table {
	display(*listdisplayA, data.FilteredList)
	return listdisplayA
}

func fillTableEnd() *tview.Table {
	display(*listdisplayE, data.SortedList)
	return listdisplayE
}

func main() {

	filter.SetDirection(tview.FlexRow).
		AddItem(untext, 0, 1, false).
		AddItem(filtersA, 0, 1, false).
		AddItem(filtersB, 0, 1, false).
		AddItem(filttext, 0, 1, false).
		AddItem(typefilters, 0, 1, false).
		AddItem(filtersZ, 0, 5, false)

	ListA.SetDirection(tview.FlexRow).
		AddItem(listinfo, 0, 1, false).
		AddItem(listdisplayA, 0, 8, true)

	ListE.SetDirection(tview.FlexRow).
		AddItem(listinfo, 0, 1, false).
		AddItem(listdisplayE, 0, 8, true)

	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 27 {
			app.Stop()
		} else if event.Rune() == 13 {
			mainForm()
			pages.SwitchToPage("Top")
		}
		return event
	})

	ListA.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 50 {
			fillTableEnd()
			pages.SwitchToPage("ListEnd")
		} else if event.Rune() == 53 {
			if data.Username != "" {
				pages.SwitchToPage("Filter")
			}
		} else if event.Rune() == 54 {
			pages.SwitchToPage("Top")
		}
		return event
	})

	ListE.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 49 {
			fillTableAlph()
			pages.SwitchToPage("ListAlph")
		} else if event.Rune() == 53 {
			pages.SwitchToPage("Filter")
		} else if event.Rune() == 54 {
			pages.SwitchToPage("Top")
		}
		return event
	})

	top.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 27 {
			pages.SwitchToPage("Top")
		}
		return event
	})

	pages.AddPage("Intro", text, true, true)
	pages.AddPage("Top", top, true, false)
	pages.AddPage("Filter", filter, true, false)
	pages.AddPage("ListAlph", ListA, true, false)
	pages.AddPage("ListEnd", ListE, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
