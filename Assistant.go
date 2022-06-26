package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	tcell "github.com/gdamore/tcell/v2"
	godotenv "github.com/joho/godotenv"
	tview "github.com/rivo/tview"
)

type Forums struct {
	Forum   Forum   `json:"data"`
	FPaging FPaging `json:"paging"`
}

type Forum struct {
	Title string `json:"title"`
	Posts []Post `json:"posts"`
}

type Post struct {
	Id         int        `json:"id"`
	Number     int        `json:"number"`
	Created_at string     `json:"created_at"`
	Created_by Created_by `json:"created_by"`
	Body       string     `json:"body"`
	Signature  string     `json:"signature"`
}

type Created_by struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Forum_avator string `json:"forum_avator"`
}

type FPaging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type PostInfo struct {
	Title string
	Name  string
}

type Entries struct {
	Entry  []Entry `json:"data"`
	Paging Paging  `json:"paging"`
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
	ForumSpoil   string
	ForumId      int
	ForumNum     int
	CompIds      []int
	CompNums     []int
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
	Dupelist     []EntryParse
	Colorlist    []ParsedEntry
	CSortedList  []ParsedEntry
}

type EntryParse struct {
	id    string
	title string
}

type ParsedEntry struct {
	id        string
	title     string
	startdate string
	enddate   string
	maltype   string
	inlist    int
}

var Postdata PostInfo

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

func Forum_Parse(spoiler string, ForumId, ForumNum int) (parselist []EntryParse) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.myanimelist.net/v2/forum/topic/"+strconv.Itoa(ForumId)+"?offset="+strconv.Itoa(ForumNum-1)+"&limit=1", nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header = http.Header{
		"X-MAL-CLIENT-ID": {SecretKey},
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

	var temp Forums

	err = json.Unmarshal(byteValue, &temp)

	if err != nil {
		fmt.Println(err)
	}
	body := temp.Forum.Posts[0].Body
	Postdata.Title = temp.Forum.Title
	Postdata.Name = temp.Forum.Posts[0].Created_by.Name
	if spoiler != "" {
		if strings.Contains(body, spoiler) {
			_, body, _ = strings.Cut(body, spoiler+"&quot;]")
			body, _, _ = strings.Cut(body, "[/spoiler]")
		}
	}
	cutbody := strings.Split(body, "url=https://myanimelist.net/anime")
	for i := 1; i < len(cutbody); i++ {
		var tempid, temptitle string
		if strings.Contains(cutbody[i], ".php?id=") {
			tempcut, _, _ := strings.Cut(cutbody[i][8:], "[/url]")
			tempid, temptitle, _ = strings.Cut(tempcut, "]")
		} else if strings.Contains(cutbody[i][1:10], "/") {
			tempcut, _, _ := strings.Cut(cutbody[i][1:], "[/url]")
			tempid, temptitle, _ = strings.Cut(tempcut, "/")
			_, temptitle, _ = strings.Cut(temptitle, "]")
		} else {
			tempcut, _, _ := strings.Cut(cutbody[i][1:], "[/url]")
			tempid, temptitle, _ = strings.Cut(tempcut, "]")
		}
		var tempEparse EntryParse
		tempEparse.id = tempid
		tempEparse.title = temptitle
		parselist = append(parselist, tempEparse)
	}
	return
}

func Comp_Parse(ForumId, ForumNum int, CompIds, CompNums []int) (complist []EntryParse) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.myanimelist.net/v2/forum/topic/"+strconv.Itoa(ForumId)+"?offset="+strconv.Itoa(ForumNum-1)+"&limit=1", nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header = http.Header{
		"X-MAL-CLIENT-ID": {SecretKey},
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

	var temp Forums

	err = json.Unmarshal(byteValue, &temp)

	if err != nil {
		fmt.Println(err)
	}
	var comp []Forums
	if len(CompIds) > 0 {
		for i := 0; i < len(CompIds); i++ {
			client := http.Client{}
			req, err := http.NewRequest("GET", "https://api.myanimelist.net/v2/forum/topic/"+strconv.Itoa(CompIds[i])+"?offset="+strconv.Itoa(CompNums[i]-1)+"&limit=1", nil)
			if err != nil {
				fmt.Println(err)
			}

			req.Header = http.Header{
				"X-MAL-CLIENT-ID": {SecretKey},
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

			err = json.Unmarshal(byteValue, &comp)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
	body := temp.Forum.Posts[0].Body
	cutbody := strings.Split(body, "url=https://myanimelist.net/anime/")
	for i := 1; i < len(cutbody); i++ {
		tempcut, _, _ := strings.Cut(cutbody[i], "[/url]")
		tempid, temptitle, _ := strings.Cut(tempcut, "]")
		var tempEparse EntryParse
		tempEparse.id = tempid
		tempEparse.title = temptitle
		complist = append(complist, tempEparse)
	}
	if len(CompIds) > 0 {
		for j := 0; j < len(CompIds); j++ {
			body := comp[j].Forum.Posts[0].Body
			cutbody := strings.Split(body, "url=https://myanimelist.net/anime/")
			for i := 1; i < len(cutbody); i++ {
				tempcut, _, _ := strings.Cut(cutbody[i], "[/url]")
				tempid, temptitle, _ := strings.Cut(tempcut, "]")
				var tempEparse EntryParse
				tempEparse.id = tempid
				tempEparse.title = temptitle
				complist = append(complist, tempEparse)
			}
		}
	}
	return
}

func Check_dupes(complist []EntryParse) (dupelist []EntryParse) {
	for i := 0; i < len(complist); i++ {
		for j := i + 1; j < len(complist); j++ {
			if complist[i].id == complist[j].id || complist[i].title == complist[j].title {
				dupelist = append(dupelist, complist[i])
			}
		}
	}
	return
}

func Check_forum(parselist []EntryParse, filtlist []CleanEntry, fulllist []Entry) (colorlist []ParsedEntry) {
	for i := 0; i < len(parselist); i++ {
		check, j := Contains(parselist[i], filtlist)
		if check {
			var temp ParsedEntry
			temp.id = filtlist[j].id
			temp.title = filtlist[j].title
			temp.startdate = filtlist[j].startdate
			temp.enddate = filtlist[j].enddate
			temp.maltype = filtlist[j].maltype
			temp.inlist = 0
			colorlist = append(colorlist, temp)
		} else {
			check, j := Contains2(parselist[i], fulllist)
			if check {
				var temp ParsedEntry
				temp.id = strconv.Itoa(fulllist[j].Node.Id)
				temp.title = fulllist[j].Node.Title
				temp.startdate = fulllist[j].List_status.Start_date
				temp.enddate = fulllist[j].List_status.Finish_date
				temp.maltype = fulllist[j].Node.Media
				temp.inlist = 1
				colorlist = append(colorlist, temp)
			} else {
				var temp ParsedEntry
				temp.id = parselist[i].id
				temp.title = parselist[i].title
				temp.inlist = 2
				colorlist = append(colorlist, temp)
			}
		}
	}
	return
}

func Contains(check EntryParse, list2 []CleanEntry) (contain bool, index int) {
	for j := 0; j < len(list2); j++ {
		if check.id == list2[j].id {
			return true, j
		}
	}
	return false, 0
}

func Contains2(check EntryParse, list2 []Entry) (contain bool, index int) {
	for j := 0; j < len(list2); j++ {
		if check.id == strconv.Itoa(list2[j].Node.Id) {
			return true, j
		}
	}
	return false, 0
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

func Date_sort2(list []ParsedEntry) (sorted_list []ParsedEntry) {
	for i := 0; i < len(list); i++ {
		sorted_list = append(sorted_list, list[i])
	}
	sort.SliceStable(sorted_list, func(i, j int) bool {
		return sorted_list[i].enddate < sorted_list[j].enddate
	})

	return
}

func Get_list(username string) (list []Entry) {

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.myanimelist.net/v2/users/"+username+"/animelist?status=completed&limit=1000&fields=list_status,media_type&nsfw=true", nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header = http.Header{
		"X-MAL-CLIENT-ID": {SecretKey},
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
			"X-MAL-CLIENT-ID": {SecretKey},
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

var SecretKey string
var Once = true
var FOnce = true
var app = tview.NewApplication()
var pages = tview.NewPages()
var top = tview.NewForm()
var filter = tview.NewFlex()
var ListA = tview.NewFlex()
var ListE = tview.NewFlex()
var ForumA = tview.NewFlex()
var ForumE = tview.NewFlex()
var Fview = tview.NewFlex()
var ForumsFlex = tview.NewFlex()
var Forums1 = tview.NewForm()
var Forums2 = tview.NewForm()
var ForumsButton = tview.NewForm()
var ForumsAdd = tview.NewForm()
var ForumsComp []*tview.Form
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
var dupedisplay = tview.NewTable().
	SetEvaluateAllRows(true).
	SetSelectable(true, false)
var forumdisplayA = tview.NewTable().
	SetEvaluateAllRows(true).
	SetSelectable(true, false)
var forumdisplayE = tview.NewTable().
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

var ForumsText1 = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("Following fields are optional, leave blank to skip")

var ForumsText2 = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("Posts to compare to:")

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

	filtersZ.AddButton("Forum Parser", func() {
		data.FilteredList = Date_filter(data.anime, data.StartYear+"-"+data.StartMonth+"-"+data.StartDay, data.EndYear+"-"+data.EndMonth+"-"+data.EndDay, data.PrevStart, data.PrevComp)
		if data.TV || data.Movie || data.OVA || data.ONA || data.Special || data.Music || data.Unknown {
			data.FilteredList = Type_filter(data.FilteredList, data.TV, data.Movie, data.OVA, data.ONA, data.Special, data.Music, data.Unknown)
		}
		data.SortedList = Date_sort(data.FilteredList)

		if FOnce {
			ForumEntry()
			ForumFilter()
			ForumButton1()
			ForumButtonZ()
		}
		FOnce = false
		pages.SwitchToPage("ForumView")
	})

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

func ForumEntry() *tview.Form {
	Forums1.SetHorizontal(true)

	Forums1.AddInputField("Forum ID", "", 8, nil, func(forumId string) {
		data.ForumId, _ = strconv.Atoi(forumId)
	})

	Forums1.AddInputField("Post Number", "", 4, nil, func(forumNum string) {
		data.ForumNum, _ = strconv.Atoi(forumNum)
	})
	return Forums1
}

func ForumFilter() *tview.Form {
	Forums2.AddInputField("Filter by Spoiler", "", 12, nil, func(forumSpoil string) {
		data.ForumSpoil = forumSpoil
	})

	return Forums2
}

func ForumButton1() *tview.Form {
	ForumsAdd.AddButton("Add", func() {
		ForumsComp = append(ForumsComp, tview.NewForm())
		ForumsFlex.AddItem(ForumsComp[len(ForumsComp)-1], 0, 1, false)
		ForumButtonN(ForumsComp[len(ForumsComp)-1])
	})

	return ForumsAdd
}

func ForumButtonN(temp *tview.Form) {
	temp.SetHorizontal(true)
	temp.AddInputField("Forum ID:", strconv.Itoa(data.ForumId), 8, nil, func(forumID2 string) {
		tempId, _ := strconv.Atoi(forumID2)
		data.CompIds = append(data.CompIds, tempId)
	})

	temp.AddInputField("Post Number:", "", 8, nil, func(forumNum2 string) {
		tempNum, _ := strconv.Atoi(forumNum2)
		data.CompIds = append(data.CompIds, tempNum)
	})

}
func ForumButtonZ() *tview.Form {
	ForumsButton.AddButton("Continue", func() {
		forumlist := Forum_Parse(data.ForumSpoil, data.ForumId, data.ForumNum)
		complist := Comp_Parse(data.ForumId, data.ForumNum, data.CompIds, data.CompNums)
		data.Dupelist = Check_dupes(complist)
		data.Colorlist = Check_forum(forumlist, data.FilteredList, data.anime)
		data.CSortedList = Date_sort2(data.Colorlist)
		listinfo.SetText("Sort: (1) Alphabetically \t (2) End Date \t\t Control: (5) Back \t (6) Start\n" + data.Username + ": " + data.StartYear + "-" + data.StartMonth + "-" + data.StartDay + " " + data.EndYear + "-" + data.EndMonth + "-" + data.EndDay + "\t\t" + Postdata.Name + ": " + Postdata.Title + "\nDuplicates:")
		forumdisplayA.Clear()
		forumdisplayE.Clear()
		dupedisplay.Clear()
		displaydupes(data.Dupelist)
		fillForumAlph()
		pages.SwitchToPage("ForumAlph")
	})
	return ForumsButton
}

func fillForumAlph() *tview.Table {
	displayForum(*forumdisplayA, data.Colorlist)
	return forumdisplayA
}

func fillForumEnd() *tview.Table {
	displayForum(*forumdisplayA, data.CSortedList)
	return forumdisplayA
}

func displayForum(forumdisplay tview.Table, colorlist []ParsedEntry) {
	c := 0
	for i := 0; i < len(colorlist); i++ {
		var color tcell.Color

		if colorlist[i].inlist == 0 {
			color = tcell.ColorGreen
		} else if colorlist[i].inlist == 1 {
			color = tcell.ColorYellow
		} else {
			color = tcell.ColorRed
		}
		var count = tview.NewTableCell("").SetTextColor(color)
		var id = tview.NewTableCell("").SetTextColor(color)
		var titles []string
		var mtype = tview.NewTableCell("").SetTextColor(color)
		var start = tview.NewTableCell("").SetTextColor(color)
		var end = tview.NewTableCell("").SetTextColor(color)

		count.SetText(strconv.Itoa(i + 1))
		id.SetText(colorlist[i].id)
		if len(colorlist[i].title) > 50 {
			temp := tview.WordWrap(colorlist[i].title, 50)
			for j := 0; j < len(temp); j++ {
				titles = append(titles, temp[j])
			}
		} else {
			titles = append(titles, colorlist[i].title)
		}
		mtype.SetText(colorlist[i].maltype)
		start.SetText(colorlist[i].startdate)
		end.SetText(colorlist[i].enddate)

		forumdisplay.SetCell(i+c, 0, count)
		forumdisplay.SetCell(i+c, 1, id)
		if len(titles) > 1 {
			for j := 0; j < len(titles); j++ {
				forumdisplay.SetCell(i+c, 2, tview.NewTableCell(titles[j]).SetTextColor(color))
				c = c + 1
			}
			c = c - 1
		} else {
			forumdisplay.SetCell(i+c, 2, tview.NewTableCell(titles[0]).SetTextColor(color))
		}
		forumdisplay.SetCell(i+c, 3, mtype)
		forumdisplay.SetCell(i+c, 4, start)
		forumdisplay.SetCell(i+c, 5, end)

	}
}

func displaydupes(dupelist []EntryParse) {
	c := 0
	for i := 0; i < len(dupelist); i++ {
		var count = tview.NewTableCell("")
		var id = tview.NewTableCell("")
		var titles []string

		count.SetText(strconv.Itoa(i + 1))
		id.SetText(dupelist[i].id)
		if len(dupelist[i].title) > 50 {
			temp := tview.WordWrap(dupelist[i].title, 50)
			for j := 0; j < len(temp); j++ {
				titles = append(titles, temp[j])
			}
		} else {
			titles = append(titles, dupelist[i].title)
		}

		dupedisplay.SetCell(i+c, 0, count)
		dupedisplay.SetCell(i+c, 1, id)
		if len(titles) > 1 {
			for j := 0; j < len(titles); j++ {
				dupedisplay.SetCell(i+c, 2, tview.NewTableCell(titles[j]))
				c = c + 1
			}
			c = c - 1
		} else {
			dupedisplay.SetCell(i+c, 2, tview.NewTableCell(titles[0]))
		}
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	SecretKey = os.Getenv("SECRET_KEY")

	filter.SetDirection(tview.FlexRow).
		AddItem(untext, 0, 1, false).
		AddItem(filtersA, 0, 1, true).
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

	ForumA.SetDirection(tview.FlexRow).
		AddItem(listinfo, 0, 1, false).
		AddItem(dupedisplay, 0, 2, true).
		AddItem(forumdisplayA, 0, 5, true)

	ForumE.SetDirection(tview.FlexRow).
		AddItem(listinfo, 0, 1, false).
		AddItem(dupedisplay, 0, 3, true).
		AddItem(forumdisplayE, 0, 5, true)

	Fview.SetDirection(tview.FlexRow).
		AddItem(Forums1, 0, 1, true).
		AddItem(ForumsText1, 0, 1, false).
		AddItem(Forums2, 0, 1, true).
		AddItem(ForumsText2, 0, 1, false).
		AddItem(ForumsFlex, 0, 5, false).
		AddItem(ForumsButton, 0, 1, false)

	ForumsFlex.SetDirection(tview.FlexRow).
		AddItem(ForumsAdd, 0, 1, true)

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

	ListA.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 50 {
			fillTableEnd()
			pages.SwitchToPage("ForumEnd")
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
		if event.Rune() == 50 {
			fillTableEnd()
			pages.SwitchToPage("ForumAlph")
		} else if event.Rune() == 53 {
			if data.Username != "" {
				pages.SwitchToPage("Filter")
			}
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
	pages.AddPage("ForumView", Fview, true, false)
	pages.AddPage("ListAlph", ListA, true, false)
	pages.AddPage("ListEnd", ListE, true, false)
	pages.AddPage("ForumAlph", ForumA, true, false)
	pages.AddPage("ForumEnd", ForumA, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
