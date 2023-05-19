package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	goquery "github.com/PuerkitoBio/goquery"
	tcell "github.com/gdamore/tcell/v2"
	colly "github.com/gocolly/colly"
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
	NumEpisodes  int          `json:"num_episodes"`
	AvgDuration  int          `json:"average_episode_duration"`
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
	id          string
	title       string
	startdate   string
	enddate     string
	maltype     string
	numepisodes int
	lengt       int
}

type ScrapedEntries struct {
	id  int
	cat string
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
	ForumTitle   string
	Slug         string
	Genre        string
	ForumId      int
	ForumNum     int
	MinLength    int
	CompIds      []int
	CompNums     []int
	HoFIds       []ScrapedEntries
	PrevStart    bool
	PrevComp     bool
	TV           bool
	OVA          bool
	ONA          bool
	Special      bool
	Movie        bool
	Music        bool
	Unknown      bool
	OrderedStart bool
	OrderedEnd   bool
	anime        []CleanEntry
	Dupelist     []EntryParse
	Colorlist    []ParsedEntry
}

type EntryParse struct {
	id     string
	title  string
	indexa int
	indexb int
}

type ParsedEntry struct {
	id          string
	title       string
	startdate   string
	enddate     string
	maltype     string
	numepisodes int
	lengt       int
	inlist      int
}

type Challenge struct {
	title   string
	id      string
	slug    string
	prevS   int
	types   int
	length  int
	ordered int
	comp    bool
}

func TurninParse(spoiler string, ForumId, ForumNum int) (parselist []EntryParse) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.myanimelist.net/v2/forum/topic/"+strconv.Itoa(ForumId)+"?offset="+strconv.Itoa(ForumNum-1)+"&limit=1", nil)
	if err != nil {
		log.Println(err)
	}

	req.Header = http.Header{
		"X-MAL-CLIENT-ID": {SecretKey},
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	byteValue, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
	}

	var temp Forums

	err = json.Unmarshal(byteValue, &temp)

	if err != nil {
		log.Println(err)
	}
	body := temp.Forum.Posts[0].Body
	data.ForumTitle = temp.Forum.Title
	data.Username = temp.Forum.Posts[0].Created_by.Name
	if spoiler != "" {
		if strings.Contains(body, spoiler) {
			_, body, _ = strings.Cut(body, spoiler+"&quot;]")
			body, _, _ = strings.Cut(body, "[/spoiler]")
		}
	}
	cutbody := strings.Split(body, "[url=")
	parselist = ASplit(cutbody, parselist)
	return
}

func CompParse(ForumId, ForumNum int, CompIds, CompNums []int) (complist []EntryParse) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.myanimelist.net/v2/forum/topic/"+strconv.Itoa(ForumId)+"?offset="+strconv.Itoa(ForumNum-1)+"&limit=1", nil)
	if err != nil {
		log.Println(err)
	}

	req.Header = http.Header{
		"X-MAL-CLIENT-ID": {SecretKey},
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	byteValue, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
	}

	var temp Forums
	var comps []Forums

	err = json.Unmarshal(byteValue, &temp)

	if err != nil {
		log.Println(err)
	}

	if len(CompIds) > 0 {
		for i := 0; i < len(CompIds); i++ {
			var comp Forums
			client := http.Client{}
			if CompIds[i] == 0 {
				break
			}
			req, err := http.NewRequest("GET", "https://api.myanimelist.net/v2/forum/topic/"+strconv.Itoa(CompIds[i])+"?offset="+strconv.Itoa(CompNums[i]-1)+"&limit=1", nil)
			if err != nil {
				log.Println(err)
			}

			req.Header = http.Header{
				"X-MAL-CLIENT-ID": {SecretKey},
			}

			resp, err := client.Do(req)

			if err != nil {
				log.Println(err)
			}

			if resp.Body != nil {
				defer resp.Body.Close()
			}

			byteValue, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				log.Println(err)
			}

			err = json.Unmarshal(byteValue, &comp)

			if err != nil {
				log.Println(err)
			}
			comps = append(comps, comp)
		}

	}
	body := temp.Forum.Posts[0].Body
	cutbody := strings.Split(body, "[url=")
	complist = ASplit(cutbody, complist)
	if len(comps) > 0 {
		for j := 0; j < len(comps); j++ {
			if len(comps[j].Forum.Posts) > 0 {
				body := comps[j].Forum.Posts[0].Body
				cutbody := strings.Split(body, "[url=")
				complist = ASplit(cutbody, complist)
			}
		}
	}
	return
}

func BlogParse(spoiler string, BlogId int) (Alist []EntryParse) {
	var Username string
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		//		fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		//		fmt.Println(r.StatusCode)
	})
	c.OnHTML(".h1", func(e *colly.HTMLElement) {
		a := e.DOM.Find("span")
		Username = strings.ReplaceAll(a.First().Text(), "'s Blog", "")
	})

	c.OnHTML(".blog_detail_content_wrapper", func(e *colly.HTMLElement) {
		l := e.DOM.Find("a")
		l.Each(func(i int, l *goquery.Selection) {
			var tempid, temptitle string
			var tempEparse EntryParse
			//			a := l.Find("a")
			href, _ := l.Attr("href")
			temptitle = l.Text()
			a, tempid, temptitle := AnimeID(href, temptitle, false)
			if a {
				tempEparse.id = tempid
				tempEparse.title = temptitle
				Alist = append(Alist, tempEparse)
			}
		})

	})

	c.Visit("https://myanimelist.net/blog.php?eid=" + strconv.Itoa(BlogId))

	data.Username = Username
	data.ForumTitle = "Blog Post"
	return
}

func BlogComps(spoiler string, BlogId []int) (Alist []EntryParse) {

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		//		fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		//		fmt.Println(r.StatusCode)
	})

	c.OnHTML(".blog_detail_content_wrapper", func(e *colly.HTMLElement) {
		l := e.DOM.Find("a")
		l.Each(func(i int, l *goquery.Selection) {
			var tempid, temptitle string
			var tempEparse EntryParse
			//			a := l.Find("a")
			href, _ := l.Attr("href")
			temptitle = l.Text()
			log.Println(href, temptitle)
			a, tempid, temptitle := AnimeID(href, temptitle, true)
			if a {
				tempEparse.id = tempid
				tempEparse.title = temptitle
				Alist = append(Alist, tempEparse)
			}

		})

	})

	for i := 0; i < len(BlogId); i++ {
		c.Visit("https://myanimelist.net/blog.php?eid=" + strconv.Itoa(BlogId[i]))
	}
	return
}

func AnimeID(url, intitle string, series bool) (good bool, id, title string) {
	good = false
	var atype string
	re := regexp.MustCompile("[0-9]+")
	if strings.Contains(url, ".net") && len(url) > 30 {
		o := strings.Index(url, ".net") + 5
		atype = url[o : o+5]
	} else {
		atype = ""
	}
	if strings.Contains(url, "myanimelist") {
		if atype == "anime" && !strings.Contains(url, "genre") && !strings.Contains(url, "producer") && !strings.Contains(url, "recommendations") && !strings.Contains(url[11:], "animelist") {
			id0 := re.FindAllString(url, -1)
			if len(id0) > 0 {
				id = id0[0]
				title = intitle
				d, b := strconv.Atoi(id)
				if series {
					if b == nil && title != "Series" && !strings.Contains(title, "http") && !(d == 199 && title == "recommended") {
						good = true
					}
				} else {
					if b == nil && !strings.Contains(title, "http") && !(d == 199 && title == "recommended") {
						good = true
					}
				}
			} else {
				log.Println("Invalid URL: ", url)
			}
		}
	}
	return
}

func ASplit(Alinks []string, Alist []EntryParse) []EntryParse {
	for i := 1; i < len(Alinks); i++ {
		var tempid, temptitle string
		var tempEparse EntryParse
		tempcut, _, _ := strings.Cut(Alinks[i], "[/url]")
		tempid, temptitle, _ = strings.Cut(tempcut, "]")
		a, tempid, temptitle := AnimeID(tempid, temptitle, true)
		if a {
			tempEparse.id = tempid
			tempEparse.title = temptitle
			Alist = append(Alist, tempEparse)
		}

	}
	return Alist
}

func CheckDupes(complist []EntryParse) (dupelist []EntryParse) {
	for i := 0; i < len(complist); i++ {
		for j := i + 1; j < len(complist); j++ {
			a, _ := strconv.Atoi(complist[i].id)
			b, _ := strconv.Atoi(complist[j].id)
			if a == b && a != 0 {
				var tmp EntryParse
				tmp.id = complist[i].id
				tmp.title = complist[i].title
				tmp.indexa = i + 1
				tmp.indexb = j + 1
				dupelist = append(dupelist, tmp)
			}
		}
	}
	return
}

func CheckForum(parselist []EntryParse, fulllist []CleanEntry) (colorlist []ParsedEntry) {
	for i := 0; i < len(parselist); i++ {
		check, j := Contains(parselist[i], fulllist)
		if check {
			var temp ParsedEntry
			temp.id = fulllist[j].id
			temp.title = fulllist[j].title
			temp.startdate = fulllist[j].startdate
			temp.enddate = fulllist[j].enddate
			temp.maltype = fulllist[j].maltype
			temp.numepisodes = fulllist[j].numepisodes
			temp.lengt = fulllist[j].lengt
			temp.inlist = 0
			colorlist = append(colorlist, temp)
		} else {
			var temp ParsedEntry
			temp.id = parselist[i].id
			temp.title = parselist[i].title
			temp.inlist = 2
			colorlist = append(colorlist, temp)

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

func Contains3(check ParsedEntry, list2 []ScrapedEntries) (contain bool, index int) {
	for j := 0; j < len(list2); j++ {
		if check.id == strconv.Itoa(list2[j].id) {
			return true, j
		}
	}
	return false, 0
}

func Get_list(username string) (CleanList []CleanEntry) {

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.myanimelist.net/v2/users/"+username+"/animelist?status=completed&limit=1000&fields=list_status,media_type,num_episodes,average_episode_duration&nsfw=true", nil)
	if err != nil {
		log.Println(err)
	}

	req.Header = http.Header{
		"X-MAL-CLIENT-ID": {SecretKey},
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	byteValue, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
	}

	var temp Entries

	err = json.Unmarshal(byteValue, &temp)

	if err != nil {
		log.Println(err)
	}

	next := temp.Paging.Next

	for next != "" {
		req, err := http.NewRequest("GET", next, nil)
		if err != nil {
			log.Println(err)
		}

		req.Header = http.Header{
			"X-MAL-CLIENT-ID": {SecretKey},
		}

		resp, err := client.Do(req)

		if err != nil {
			log.Println(err)
		}

		if resp.Body != nil {
			defer resp.Body.Close()
		}

		byteValue, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Println(err)
		}

		var temp2 Entries

		err = json.Unmarshal(byteValue, &temp2)

		if err != nil {
			log.Println(err)
		}

		for i := 0; i < len(temp2.Entry); i++ {
			temp.Entry = append(temp.Entry, temp2.Entry[i])
		}

		next = temp2.Paging.Next
	}

	list := temp.Entry

	temp3 := CleanEntry{}
	for i := 0; i < len(list); i++ {
		temp3.id = strconv.Itoa(list[i].Node.Id)
		temp3.title = list[i].Node.Title
		temp3.startdate = list[i].List_status.Start_date
		temp3.enddate = list[i].List_status.Finish_date
		temp3.maltype = list[i].Node.Media
		temp3.numepisodes = list[i].Node.NumEpisodes
		temp3.lengt = list[i].Node.NumEpisodes * list[i].Node.AvgDuration
		CleanList = append(CleanList, temp3)
	}
	return
}

func ScrapeHOF(username, slug string) (ids []ScrapedEntries) {
	t := ""
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(r.StatusCode)
	})
	c.OnHTML("#challengeItems", func(e *colly.HTMLElement) {

		r := e.DOM.Find(".listItem")
		r.Each(func(i int, r *goquery.Selection) {
			var temp = ScrapedEntries{}
			a, b := r.Attr("category")
			if b {
				//				if a != "Completed" && a != "Previously Completed" && a != "On Hold" && a != "Plan to Watch" && a != "Unwatched" {
				temp.cat = a
				//				}
			}
			if slug == "mascot" {
				a := r.Find(".entry-comments")
				temp.cat = a.Text()
			}
			e := r.Find(".seriesLink")
			a, b = e.Attr("seriesid")
			if b {
				a2, _ := strconv.Atoi(a)
				temp.id = a2
				ids = append(ids, temp)
			}
		})

	})
	c.Visit("https://anime.jhiday.net/hof/challenge/" + slug + "?user=" + username + "#challengeItems" + t)

	return
}

func ScrapeHOFAnime(id, genre string) (d string) {

	d = "removed"
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(r.StatusCode)
	})
	c.OnHTML(".table-bordered", func(e *colly.HTMLElement) {

		r := e.DOM.Find(".history-item")
		r.Each(func(i int, r *goquery.Selection) {
			a := r.Find("td")
			if strings.Contains(a.Last().Text(), "removed") {
				b := strings.Split(a.Last().Text(), "removed")[1]
				if strings.Contains(b, genre) {
					d = d + " " + a.First().Text()[:10]
				}
			}
		})
		if d == "removed" {
			d = "Not in HOF"
		}
	})
	c.Visit("https://anime.jhiday.net/hof/anime/" + id)

	return
}

var SecretKey string
var Once = true
var FOnce = true
var app = tview.NewApplication()
var pages = tview.NewPages()
var ForumInfo = tview.NewFlex()
var ForumsFlex = tview.NewFlex()
var Forums1 = tview.NewForm()
var Forums2 = tview.NewForm()
var ForumsButton = tview.NewForm()
var ForumsAdd = tview.NewForm()
var ForumsComp []*tview.Form
var CompBoxes []*tview.InputField
var filter = tview.NewFlex()
var filtersA = tview.NewForm()
var filtersB = tview.NewForm()
var filtersZ = tview.NewForm()
var typefilters = tview.NewForm()
var miscfilters = tview.NewForm()
var Months = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
var genres = map[interface{}]string{
	"action":        "Action",
	"adventure":     "Adventure",
	"comedy":        "Comedy",
	"dementia":      "Avant Garde",
	"demons":        "Mythology",
	"drama":         "Drama",
	"ecchi":         "Ecchi",
	"fantasy":       "Fantasy",
	"game":          "Game",
	"harem":         "Harem",
	"hentai":        "Hentai",
	"historical":    "Historical",
	"horror":        "Horror",
	"josei":         "Josei",
	"martialarts":   "Martial Arts",
	"mecha":         "Mecha",
	"military":      "Military",
	"musical":       "Music",
	"mystery":       "Mystery",
	"parody":        "Parody",
	"policecars":    "Detective",
	"psychological": "Psychological",
	"romance":       "Romance",
	"samurai":       "Samurai",
	"school":        "School",
	"scifi":         "Sci-Fi",
	"shoujoai":      "Girls Love",
	"shounenai":     "Boys Love",
	"sliceoflife":   "Slice of Life",
	"space":         "Space",
	"sportsv2":      "Sports",
	"superpower":    "Super Power",
	"supernatural":  "Supernatural",
	"thriller":      "Suspense",
	"vampire":       "Vampire",
}
var ForumPage = tview.NewFlex()
var dupedisplay = tview.NewTable().
	SetEvaluateAllRows(true).
	SetSelectable(true, false)
var ForumDisplay = tview.NewTable().
	SetEvaluateAllRows(true).
	SetSelectable(true, false)

var data = Filter_Data{}

var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("Press (Enter) to continue")

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

var ForumsText1 = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("Following fields are optional, leave blank to skip")

func ForumFilter() *tview.Form {
	Forums2.AddInputField("Filter by Spoiler", "", 12, nil, func(forumSpoil string) {
		data.ForumSpoil = forumSpoil
	})

	return Forums2
}

var ForumsText2 = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("Posts to compare to:")

func ForumButton1() *tview.Form {
	ForumsAdd.AddButton("Add", func() {
		data.CompIds = append(data.CompIds, 0)
		data.CompNums = append(data.CompNums, 0)
		ForumsComp = append(ForumsComp, tview.NewForm())
		ForumsFlex.AddItem(ForumsComp[len(ForumsComp)-1], 0, 1, false)
		ForumButtonN(ForumsComp[len(ForumsComp)-1])
	})

	return ForumsAdd
}

func ForumButtonN(temp *tview.Form) {
	temp.SetHorizontal(true)

	CompBoxes = append(CompBoxes, tview.NewInputField().
		SetLabel("Forum ID").
		SetFieldWidth(8))

	temp.AddFormItem(CompBoxes[len(CompBoxes)-1])

	CompBoxes = append(CompBoxes, tview.NewInputField().
		SetLabel("Post Number").
		SetFieldWidth(4))

	temp.AddFormItem(CompBoxes[len(CompBoxes)-1])

}

var listinfo = tview.NewTextView().
	SetTextColor(tcell.ColorGreen)

var untext = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetTextAlign(1)

func ForumButtonZ() *tview.Form {
	ForumsButton.AddButton("Continue", func() {
		for i := 0; i < len(ForumsComp); i++ {
			if CompBoxes[2*i].GetText() == "" {
				data.CompIds[i] = data.ForumId
			} else {
				data.CompIds[i], _ = strconv.Atoi(CompBoxes[2*i].GetText())
			}
			data.CompNums[i], _ = strconv.Atoi(CompBoxes[2*i+1].GetText())
		}
		var forumlist []EntryParse
		var complist []EntryParse
		if data.ForumNum == 0 {
			forumlist = BlogParse(data.ForumSpoil, data.ForumId)
			complist = BlogComps(data.ForumSpoil, append(data.CompIds, data.ForumId))
		} else {
			forumlist = TurninParse(data.ForumSpoil, data.ForumId, data.ForumNum)
			complist = CompParse(data.ForumId, data.ForumNum, data.CompIds, data.CompNums)
		}
		data.Dupelist = CheckDupes(complist)
		data.anime = Get_list(data.Username)
		data.Colorlist = CheckForum(forumlist, data.anime)
		untext.SetText(data.Username + "'s List")
		if Once {
			filterstart()
			filterend()
			typebuttons()
			miscbuttons()
			filterbuttons()
		}
		Once = false
		pages.SwitchToPage("Filter")
	})
	return ForumsButton
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

var filttext = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("Filter by type (leave blank to skip)")

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

func miscbuttons() *tview.Form {
	miscfilters.SetHorizontal(true)

	miscfilters.AddInputField("Min Length", "", 2, nil, func(lengt string) {
		if lengt == "" {
			data.MinLength = 0
		} else {
			a, _ := strconv.Atoi(lengt)
			data.MinLength = a * 60
		}
	})

	miscfilters.AddInputField("HoF Slug", "", 5, nil, func(slug string) {
		data.Slug = slug
	})

	miscfilters.AddInputField("Removed Genre", "", 5, nil, func(genre string) {
		data.Genre = genre
	})

	miscfilters.AddCheckbox("Ordered Start Dates", false, func(OStart bool) {
		data.OrderedStart = OStart
	})

	miscfilters.AddCheckbox("Ordered End Dates", false, func(OEnd bool) {
		data.OrderedEnd = OEnd
	})

	return miscfilters
}

func displaydupes(dupelist []EntryParse) {
	c := 0
	for i := 0; i < len(dupelist); i++ {
		var count = tview.NewTableCell("")
		var id = tview.NewTableCell("")
		var titles []string
		var indexa = tview.NewTableCell("")
		var indexb = tview.NewTableCell("")

		count.SetText(strconv.Itoa(i + 1))
		id.SetText(dupelist[i].id)
		indexa.SetText(strconv.Itoa(dupelist[i].indexb))
		indexb.SetText(strconv.Itoa(dupelist[i].indexa))
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
		dupedisplay.SetCell(i+c, 3, indexb)
		dupedisplay.SetCell(i+c, 4, indexa)

	}
}

func displayForum(colorlist []ParsedEntry) {
	c := 0
	prevEnd := "0000-00-00"

	data.HoFIds = ScrapeHOF(data.Username, data.Slug)
	for i := 0; i < len(colorlist); i++ {
		var color tcell.Color

		startday := data.StartYear + "-" + data.StartMonth + "-" + data.StartDay
		endday := data.EndYear + "-" + data.EndMonth + "-" + data.EndDay
		inlist := colorlist[i].inlist

		if data.PrevComp {
			if colorlist[i].enddate > endday {
				inlist = 1
			}
		} else if data.PrevStart {
			if colorlist[i].enddate < startday || colorlist[i].enddate > endday {
				inlist = 1
			}
		} else {
			if colorlist[i].enddate > endday || colorlist[i].startdate < startday {
				inlist = 1
			}
		}

		if inlist == 0 {
			color = tcell.ColorGreen
		} else if inlist == 1 {
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
		var hof = tview.NewTableCell("").SetTextColor(tcell.ColorGreen)
		var leng = tview.NewTableCell("").SetTextColor(color)

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
		if data.TV || data.Movie || data.OVA || data.ONA || data.Special || data.Music || data.Unknown {
			var media_types = []string{}
			var in_types = false
			if data.TV {
				media_types = append(media_types, "tv")
			}
			if data.Movie {
				media_types = append(media_types, "movie")
			}
			if data.OVA {
				media_types = append(media_types, "ova")
			}
			if data.ONA {
				media_types = append(media_types, "ona")
			}
			if data.Special {
				media_types = append(media_types, "special")
			}
			if data.Music {
				media_types = append(media_types, "music")
			}
			if data.Unknown {
				media_types = append(media_types, "unknown")
			}
			for _, val := range media_types {
				if val == colorlist[i].maltype {
					in_types = true
				}
			}
			if !in_types {
				mtype.SetTextColor(tcell.ColorRed)
			}
		}
		mtype.SetText(colorlist[i].maltype)
		start.SetText(colorlist[i].startdate)
		end.SetText(colorlist[i].enddate)
		if data.MinLength == 3*60 {
			if colorlist[i].maltype != "tv" && colorlist[i].numepisodes < 10 && (colorlist[i].numepisodes < 4 || colorlist[i].lengt < 80*60) {
				leng.SetText("Too Short ")
				leng.SetTextColor(tcell.ColorRed)
			}
		} else if data.MinLength == 2*60 {
			if colorlist[i].numepisodes < 10 && (colorlist[i].numepisodes < 4 || colorlist[i].lengt < 80*60) {
				leng.SetText("Too Short ")
				leng.SetTextColor(tcell.ColorRed)
			}
		} else if data.MinLength > 0 {
			if colorlist[i].lengt >= data.MinLength {
				leng.SetText("")
			} else {
				leng.SetText("Too Short ")
				leng.SetTextColor(tcell.ColorRed)
			}
		} else {
			leng.SetText("")
		}
		if data.Slug != "" {
			a, b := Contains3(colorlist[i], data.HoFIds)
			if !a {
				hof.SetTextColor(tcell.ColorRed)
				if data.Genre != "" {
					hof.SetText(ScrapeHOFAnime(colorlist[i].id, data.Genre))
				} else if genres[data.Slug] != "" {
					hof.SetText(ScrapeHOFAnime(colorlist[i].id, genres[data.Slug]))
				} else {
					hof.SetText("Not in HoF")
				}
			} else {
				e := data.HoFIds[b].cat
				e = strings.Replace(e, "[", "", -1)
				e = strings.Replace(e, "]", "", -1)
				hof.SetText(e)
			}
		}

		if data.OrderedStart {
			if colorlist[i].startdate < prevEnd {
				start.SetTextColor(tcell.ColorRed)
			}
		}
		if data.OrderedEnd {
			if colorlist[i].enddate < prevEnd {
				end.SetTextColor(tcell.ColorRed)
			}
		}
		prevEnd = colorlist[i].enddate

		ForumDisplay.SetCell(i+c, 0, count)
		ForumDisplay.SetCell(i+c, 1, id)
		if len(titles) > 1 {
			for j := 0; j < len(titles); j++ {
				ForumDisplay.SetCell(i+c, 2, tview.NewTableCell(titles[j]).SetTextColor(color))
				c = c + 1
			}
			c = c - 1
		} else {
			ForumDisplay.SetCell(i+c, 2, tview.NewTableCell(titles[0]).SetTextColor(color))
		}
		ForumDisplay.SetCell(i+c, 3, mtype)

		ForumDisplay.SetCell(i+c, 4, start)
		ForumDisplay.SetCell(i+c, 5, end)
		ForumDisplay.SetCell(i+c, 6, hof)
		ForumDisplay.SetCell(i+c, 7, leng)

	}
}

func filterbuttons() *tview.Form {
	filtersZ.SetHorizontal(true)

	filtersZ.AddButton("Continue", func() {

		listinfo.SetText("Control: (5) Back \t (6) Start\n" + data.Username + ": " + data.StartYear + "-" + data.StartMonth + "-" + data.StartDay + " " + data.EndYear + "-" + data.EndMonth + "-" + data.EndDay + "\t\t" + data.Username + ": " + data.ForumTitle + " " + data.ForumSpoil + "\nDuplicates:")
		ForumDisplay.Clear()
		dupedisplay.Clear()
		displaydupes(data.Dupelist)
		displayForum(data.Colorlist)
		pages.SwitchToPage("ForumView")
	})

	return filtersZ
}

func main() {

	os.Remove("logs.txt")

	var file, err = os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.Lshortfile)

	log.SetOutput(file)

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	SecretKey = os.Getenv("SECRET_KEY")

	filter.SetDirection(tview.FlexRow).
		AddItem(untext, 0, 1, false).
		AddItem(filtersA, 0, 1, true).
		AddItem(filtersB, 0, 1, false).
		AddItem(filttext, 0, 1, false).
		AddItem(typefilters, 0, 1, false).
		AddItem(miscfilters, 0, 1, false).
		AddItem(filtersZ, 0, 4, false)

	ForumPage.SetDirection(tview.FlexRow).
		AddItem(listinfo, 0, 1, false).
		AddItem(dupedisplay, 0, 2, true).
		AddItem(ForumDisplay, 0, 5, true)

	ForumInfo.SetDirection(tview.FlexRow).
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
			if FOnce {
				ForumEntry()
				ForumFilter()
				ForumButton1()
				ForumButtonZ()
			}
			FOnce = false
			pages.SwitchToPage("FInfo")
		}
		return event
	})

	ForumPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 54 {
			pages.SwitchToPage("FInfo")
		} else if event.Rune() == 53 {
			pages.SwitchToPage("Filter")
		}
		return event
	})

	pages.AddPage("Intro", text, true, true)
	pages.AddPage("FInfo", ForumInfo, true, false)
	pages.AddPage("Filter", filter, true, false)
	pages.AddPage("ForumView", ForumPage, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Println(err)
	}

}
